package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/Manolo-Esc/gommence/src/pkg/netw"
	"github.com/go-chi/chi/v5"
)

type EnrollmentHandler struct {
	service ports.EnrollmentService
	logger  logger.LoggerService
}

func NewEnrollmentHandler(service ports.EnrollmentService, logger logger.LoggerService) *EnrollmentHandler {
	return &EnrollmentHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Get all enrollments for a user
// @Description Get all enrollments for a user
// @Tags Enrollments
// @Produce json
// @Param 	userId path string true  "ID del usuario"
// @Success 200 {object} dtos.Enrollment
// @Failure 400 "Invalid data"
// @Failure 500 "Error generating response or token"
// @Router /enrollment/user/{userId} [get]
func (h *EnrollmentHandler) GetUserEnrollments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	byUser := netw.JwtGetUserInToken(ctx)
	id := chi.URLParam(r, "userId")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // Establecer un timeout para la operación
	defer cancel()

	response, errLogin := h.service.GetEnrollmentsOfUser(ctx, byUser, id)
	if errLogin != nil {
		http.Error(w, errLogin.Error(), errLogin.Status())
		return
	}

	if err := netw.Encode(w, r, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
