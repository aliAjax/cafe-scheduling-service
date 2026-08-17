package user

import (
	"net/http"
	"strings"

	"cafe-scheduling-api/pkg/apperror"
	"cafe-scheduling-api/pkg/auth"
	"cafe-scheduling-api/pkg/httputil"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "username and password are required")
		return
	}

	result, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, result)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "authentication required")
		return
	}

	u, err := h.service.GetByID(r.Context(), claims.UserID)
	if err != nil {
		httputil.Error(w, apperror.Status(err), err.Error())
		return
	}

	httputil.JSON(w, http.StatusOK, u)
}
