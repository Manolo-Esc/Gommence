package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/dtos"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/Manolo-Esc/gommence/src/pkg/netw"
)

type AuthHandler struct {
	service ports.AuthService
	logger  logger.LoggerService
}

func NewAuthHandler(service ports.AuthService, logger logger.LoggerService) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// @Summary Sign in the system
// @Description Receives login credentials and returns a token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   loginData  body dtos.LoginCredentials  true  "Credentials"
// @Success 200 {object} dtos.LoggedUser
// @Failure 400 "Invalid data"
// @Failure 500 "Error generating response or token"
// @Router /auth/signin [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	credentials, err := netw.Decode[dtos.LoginCredentials](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	response, errLogin := h.service.Login(ctx, credentials)
	if errLogin != nil {
		http.Error(w, errLogin.Error(), errLogin.Status())
		return
	}

	if err = netw.Encode(w, r, http.StatusOK, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
