package subscription

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	respond "github.com/estenity95/go-test-task/api/resource/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type API struct {
	logger     *zerolog.Logger
	validator  *validator.Validate
	repository *Repo
}

func New(logger *zerolog.Logger, validator *validator.Validate, db *gorm.DB) *API {
	return &API{
		logger:     logger,
		validator:  validator,
		repository: NewRepo(db),
	}
}

// List обрабатывает запрос на получение списка подписок.
// Принимает необязательные query-параметры:
//   - limit  (int, по умолчанию 50, верхняя граница 1000)
//   - offset (int, по умолчанию 0, неотрицательный)
//
// Возвращает JSON-массив DTO (Resp).
//
// List godoc
// @summary     List subscriptions
// @description Returns a paginated list of subscriptions
// @tags        subscriptions
// @produce     json
// @param       limit   query int false "items per page" default(50)
// @param       offset  query int false "offset"         default(0)
// @success     200 {array}  subscription.Resp
// @failure     500 {object} respond.Problem
// @router      /subscriptions [get]
func (h *API) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)

	q := r.URL.Query()

	limit := atoiDefault(q.Get("limit"), 50)
	offset := atoiDefault(q.Get("offset"), 0)

	if limit < 1 {
		limit = 1
	} else if limit > 1000 {
		limit = 1000
	}

	if offset < 0 {
		offset = 0
	}

	items, err := h.repository.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error().Str("rid", reqID).Err(err).Msg("list subscriptions failed")
		respond.Internal(w, "failed to list subscriptions")
		return
	}

	out := make([]Resp, 0, len(items)) // гарантируем [] а не null
	for _, s := range items {
		out = append(out, toResp(*s))
	}

	respond.OK(w, out)
}

// Create создаёт новую подписку.
//
// Create godoc
// @summary     Create subscription
// @tags        subscriptions
// @accept      json
// @produce     json
// @param       body body      subscription.CreateReq true "Create payload"
// @success     201  {object}  subscription.Resp
// @failure     400  {object}  respond.Problem
// @failure     500  {object}  respond.Problem
// @router      /subscriptions [post]
func (h *API) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)

	var in CreateReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.logger.Error().Str("rid", reqID).Err(err).Msg("decode create body failed")
		respond.BadRequest(w, "invalid json body")
		return
	}

	if err := h.validator.Struct(in); err != nil {
		respond.ValidationFailed(w, mapValidationErrors(err))
		return
	}

	ent, err := toEntity(in)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	created, err := h.repository.Create(ctx, &ent)
	if err != nil {
		h.logger.Error().Str("rid", reqID).Err(err).Msg("create subscription failed")
		respond.Internal(w, "failed to create subscription")
		return
	}

	respond.Created(w, toResp(*created))
}

// Read возвращает одну подписку по идентификатору.
//
// Read godoc
// @summary     Read subscription
// @tags        subscriptions
// @produce     json
// @param       id   path int true "ID"
// @success     200  {object}  subscription.Resp
// @failure     400  {object}  respond.Problem
// @failure     404  {object}  respond.Problem
// @failure     500  {object}  respond.Problem
// @router      /subscriptions/{id} [get]
func (h *API) Read(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respond.BadRequest(w, "id must be integer")
		return
	}

	item, err := h.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(w, "subscription not found")
			return
		}
		h.logger.Error().Str("rid", reqID).Err(err).Msg("read subscription failed")
		respond.Internal(w, "failed to read subscription")
		return
	}

	respond.OK(w, toResp(*item))
}

// Update обновляет подписку целиком.
//
// Update godoc
// @summary     Update subscription
// @tags        subscriptions
// @accept      json
// @produce     json
// @param       id   path int                    true "ID"
// @param       body body subscription.UpdateReq true "Update payload"
// @success     200  {object}  subscription.Resp
// @failure     400  {object}  respond.Problem
// @failure     404  {object}  respond.Problem
// @failure     500  {object}  respond.Problem
// @router      /subscriptions/{id} [put]
func (h *API) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respond.BadRequest(w, "id must be integer")
		return
	}

	var in UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		h.logger.Error().Str("rid", reqID).Err(err).Msg("decode update body failed")
		respond.BadRequest(w, "invalid json body")
		return
	}

	if err := h.validator.Struct(in); err != nil {
		respond.ValidationFailed(w, mapValidationErrors(err))
		return
	}

	ent, err := toEntity(in)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	ent.ID = id

	_, err = h.repository.Update(ctx, &ent)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(w, "subscription not found")
			return
		}
		h.logger.Error().Str("rid", reqID).Err(err).Msg("update subscription failed")
		respond.Internal(w, "failed to update subscription")
		return
	}

	got, err := h.repository.Get(ctx, id)
	if err != nil {
		h.logger.Error().Str("rid", reqID).Err(err).Msg("fetch after update failed")
		respond.Internal(w, "failed to fetch updated subscription")
		return
	}

	respond.OK(w, toResp(*got))
}

// Delete удаляет подписку.
//
// Delete godoc
// @summary     Delete subscription
// @tags        subscriptions
// @param       id   path int true "ID"
// @success     204  {string} string "no content"
// @failure     400  {object} respond.Problem
// @failure     404  {object} respond.Problem
// @failure     500  {object} respond.Problem
// @router      /subscriptions/{id} [delete]
func (h *API) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respond.BadRequest(w, "id must be integer")
		return
	}

	if err := h.repository.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			respond.NotFound(w, "subscription not found")
			return
		}
		h.logger.Error().Str("rid", reqID).Err(err).Msg("delete subscription failed")
		respond.Internal(w, "failed to delete subscription")
		return
	}

	respond.NoContent(w)
}

// Summary возвращает сумму price всех подписок, которые пересекаются с заданным периодом.
// Фильтры: user_id (UUID), service_name (string). Период и фильтры — через query.
//
// Поддерживаемые форматы дат в query:
//   - YYYY-MM:      2025-01  (будет интерпретирован как первое число месяца)
//
// Пример:
//
//	GET /api/v1/subscriptions/summary?from=2025-01&to=2025-06&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Netflix
//
// Summary godoc
// @summary     Sum subscription prices over a period
// @description Sums `price` for all subscriptions overlapping [from..to]; optional filters by user_id and service_name
// @tags        subscriptions
// @produce     json
// @param       from         query string true  "start of period (RFC3339 | 2006-01-02 | 2006-01)"
// @param       to           query string true  "end of period   (RFC3339 | 2006-01-02 | 2006-01)"
// @param       user_id      query string false "user UUID filter"
// @param       service_name query string false "service name filter"
// @success     200 {object} map[string]int64 "example: {\"total\": 2397}"
// @failure     400 {object} respond.Problem
// @failure     500 {object} respond.Problem
// @router      /subscriptions/summary [get]
func (h *API) Summary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := middleware.GetReqID(ctx)

	q := r.URL.Query()
	fromStr := q.Get("from")
	toStr := q.Get("to")
	if fromStr == "" || toStr == "" {
		respond.BadRequest(w, "`from` and `to` are required")
		return
	}

	from, err := parseDateQuery(fromStr)
	if err != nil {
		respond.BadRequest(w, "invalid `from` date")
		return
	}
	to, err := parseDateQuery(toStr)
	if err != nil {
		respond.BadRequest(w, "invalid `to` date")
		return
	}
	if to.Before(from) {
		respond.BadRequest(w, "`to` must be >= `from`")
		return
	}

	// опциональные фильтры
	var uidPtr *uuid.UUID
	if uid := q.Get("user_id"); uid != "" {
		u, err := uuid.Parse(uid)
		if err != nil {
			respond.BadRequest(w, "`user_id` must be a valid UUID")
			return
		}
		uidPtr = &u
	}
	var svcPtr *string
	if name := q.Get("service_name"); name != "" {
		svc := name
		svcPtr = &svc
	}

	total, err := h.repository.Summary(ctx, SummaryQuery{
		From:        from,
		To:          to,
		UserID:      uidPtr,
		ServiceName: svcPtr,
	})

	if err != nil {
		h.logger.Error().Str("rid", reqID).Err(err).Msg("summary failed")
		respond.Internal(w, "failed to compute summary")
		return
	}

	respond.OK(w, map[string]int64{"total": total})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}

	return n
}

func mapValidationErrors(err error) map[string]string {
	out := map[string]string{}
	if err == nil {
		return out
	}
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fe := range validationErrors {
			out[fe.Field()] = fe.Tag()
		}

		return out
	}
	out["error"] = err.Error()

	return out
}

func parseDateQuery(s string) (time.Time, error) {
	if t, err := time.Parse("01-2006", s); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("bad date: %q", s)
}
