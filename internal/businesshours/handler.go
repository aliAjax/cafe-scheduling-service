package businesshours

import (
	"net/http"

	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, items)
}

func (h *Handler) UpsertWeekly(w http.ResponseWriter, r *http.Request) {
	var req []UpsertItem
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON body; expected an array")
		return
	}
	items, err := h.service.UpsertWeekly(r.Context(), req)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, items)
}
