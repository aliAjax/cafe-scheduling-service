package schedule

import (
	"net/http"
	"strconv"
	"time"

	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/auth"
	"cafe-scheduling-api/pkg/httputil"
	"cafe-scheduling-api/pkg/timeutil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	weekStart := r.URL.Query().Get("week_start")
	if weekStart == "" {
		weekStart = timeutil.FormatDate(timeutil.StartOfWeek(time.Now()))
	}

	var employeeID *int64
	if claims.Role == "employee" {
		if claims.EmployeeID == nil {
			httputil.Error(w, http.StatusForbidden, "employee account is not linked to an employee profile")
			return
		}
		employeeID = claims.EmployeeID
	} else if raw := r.URL.Query().Get("employee_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, "invalid employee_id query parameter")
			return
		}
		employeeID = &id
	}

	items, err := h.service.List(r.Context(), weekStart, employeeID)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, items)
}

func (h *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	if claims.EmployeeID == nil {
		httputil.Error(w, http.StatusForbidden, "employee account is not linked to an employee profile")
		return
	}
	weekStart := r.URL.Query().Get("week_start")
	if weekStart == "" {
		weekStart = timeutil.FormatDate(timeutil.StartOfWeek(time.Now()))
	}
	items, err := h.service.List(r.Context(), weekStart, claims.EmployeeID)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, items)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	item, err := h.service.Create(r.Context(), req)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusCreated, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid schedule id")
		return
	}
	var req CreateRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	item, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, item)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid schedule id")
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, httputil.MessageResponse{Message: "schedule deleted"})
}
