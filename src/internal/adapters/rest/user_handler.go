package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/Manolo-Esc/gommence/src/pkg/netw"
)

type UserHandler struct {
	service ports.UserService
	logger  logger.LoggerService
}

func NewUserHandler(service ports.UserService, logger logger.LoggerService) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Get all Users
// @Description Get all Users in the system
// @Tags Users
// @Produce json
// @Success 200 {array} dtos.User
// @Failure 400 "Invalid data"
// @Failure 500 "Error generating response"
// @Router /user/ [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	byUser := netw.JwtGetUserInToken(ctx)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // Establecer un timeout para la operación
	defer cancel()

	response, errLogin := h.service.GetUsers(ctx, byUser)
	if errLogin != nil {
		http.Error(w, errLogin.Error(), errLogin.Status())
		return
	}

	if err := netw.Encode(w, r, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get all Users
// @Description Get all Users in the system
// @Tags Users
// @Produce json
// @Param 	userId path string true  "ID del usuario"
// @Success 200 {object} dtos.User
// @Failure 400 "Invalid data"
// @Failure 500 "Error generating response or token"
// @Router /User/user/{userId} [get]
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// byUser := netw.JwtGetUserInToken(ctx)
	// id := chi.URLParam(r, "userId")

	// ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // Establecer un timeout para la operación
	// defer cancel()

	// response, errLogin := h.service.GetUsersOfUser(ctx, byUser, id)
	// if errLogin != nil {
	// 	http.Error(w, errLogin.Error(), errLogin.Status())
	// 	return
	// }

	// if err := netw.Encode(w, r, http.StatusOK, response); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}
