package subscription

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64      `gorm:"primaryKey;column:id"`
	ServiceName string     `gorm:"column:service_name;not null"`
	Price       int        `gorm:"column:price;not null"`
	UserID      uuid.UUID  `gorm:"type:uuid;column:user_id;not null"`
	StartMonth  time.Time  `gorm:"column:start_month;not null"`
	EndMonth    *time.Time `gorm:"column:end_month"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

type Subscriptions []*Subscription

type CreateReq struct {
	ServiceName string  `json:"service_name"  validate:"required"`
	Price       int     `json:"price"         validate:"gte=0"`
	UserID      string  `json:"user_id"       validate:"required,uuid4"`
	StartDate   string  `json:"start_date"    validate:"required"`
	EndDate     *string `json:"end_date"      validate:"omitempty"`
}

type UpdateReq = CreateReq

type Resp struct {
	ID          int64   `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

var (
	ErrBadInput   = errors.New("bad input")
	ErrNotFound   = errors.New("not found")
	ErrBadDate    = errors.New("bad date")
	timeLayoutOut = "2006-01-02T15:04:05Z07:00"
)

func parseDate(s string) (time.Time, error) {
	if t, err := time.Parse("01-2006", s); err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}

	return time.Time{}, ErrBadDate
}

func toEntity(in CreateReq) (Subscription, error) {
	uid, err := uuid.Parse(in.UserID)
	if err != nil {
		return Subscription{}, err
	}

	sm, err := parseDate(in.StartDate)
	if err != nil {
		return Subscription{}, err
	}

	var em *time.Time
	if in.EndDate != nil && *in.EndDate != "" {
		t, err := parseDate(*in.EndDate)
		if err != nil {
			return Subscription{}, err
		}
		em = &t
	}

	if em != nil && em.Before(sm) {
		return Subscription{}, ErrBadInput
	}

	return Subscription{
		ServiceName: in.ServiceName,
		Price:       in.Price,
		UserID:      uid,
		StartMonth:  sm,
		EndMonth:    em,
	}, nil
}

func toResp(s Subscription) Resp {
	out := Resp{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID.String(),
		StartDate:   s.StartMonth.Format("2006-01-02"),
		CreatedAt:   s.CreatedAt.Format(timeLayoutOut),
		UpdatedAt:   s.UpdatedAt.Format(timeLayoutOut),
	}

	if s.EndMonth != nil {
		e := s.EndMonth.Format("2006-01-02")
		out.EndDate = &e
	}
	return out
}
