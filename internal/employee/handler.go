package employee

import (
	"net/http"
	"strconv"

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
	employees, err := h.service.List(r.Context())
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, employees)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	employee, err := h.service.Create(r.Context(), req)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusCreated, employee)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid employee id")
		return
	}
	employee, err := h.service.Get(r.Context(), id)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, employee)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	var req UpdateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	employee, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, employee)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid employee id")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, httputil.MessageResponse{Message: "employee deleted"})
}
