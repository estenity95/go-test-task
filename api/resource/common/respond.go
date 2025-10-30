package respond

import (
	"encoding/json"
	"net/http"
)

type Problem struct {
	Type     string            `json:"type,omitempty"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail,omitempty"`
	Instance string            `json:"instance,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

type envelope[T any] struct {
	Data T `json:"data"`
}

func toJson(w http.ResponseWriter, code int, v any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func OK[T any](w http.ResponseWriter, v T) {
	toJson(w, http.StatusOK, envelope[T]{Data: v})
}

func Created[T any](w http.ResponseWriter, v T) {
	toJson(w, http.StatusCreated, envelope[T]{Data: v})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, msg string) {
	toJson(w, http.StatusBadRequest, Problem{Title: "bad request", Status: 400, Detail: msg})
}
func NotFound(w http.ResponseWriter, msg string) {
	toJson(w, http.StatusNotFound, Problem{Title: "not found", Status: 404, Detail: msg})
}

func Internal(w http.ResponseWriter, msg string) {
	toJson(w, http.StatusInternalServerError, Problem{Title: "internal error", Status: 500, Detail: msg})
}

func ValidationFailed(w http.ResponseWriter, details map[string]string) {
	toJson(w, http.StatusBadRequest, Problem{
		Title:  "validation failed",
		Status: http.StatusBadRequest,
		Errors: details,
	})
}
