package statistics

import (
	"net/http"
	"time"

	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/httputil"
	"cafe-scheduling-api/pkg/timeutil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Weekly(w http.ResponseWriter, r *http.Request) {
	weekStart := r.URL.Query().Get("week_start")
	if weekStart == "" {
		weekStart = timeutil.FormatDate(timeutil.StartOfWeek(time.Now()))
	}
	stats, err := h.service.Weekly(r.Context(), weekStart)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, stats)
}
