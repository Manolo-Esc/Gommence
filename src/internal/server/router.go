package server

import (
	"net/http"

	_ "github.com/Manolo-Esc/gommence/src/docs"
	"github.com/Manolo-Esc/gommence/src/internal/adapters/rest"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/Manolo-Esc/gommence/src/pkg/netw"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// These are annotations for Swagger documentation
// @title           Gommence
// @version         1.0
// @description     Go Web Server starter kit
// @host           localhost:8080
// @BasePath       /api/v1

// @Summary Health checking URL
// @Tags Misc
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"mensaje": "service is online"}`))
}

func addRoutes(appModules *AppModules, r *chi.Mux, logger logger.LoggerService, db *gorm.DB) {
	authHandler := rest.NewAuthHandler(*appModules.auth, logger)
	userHandler := rest.NewUserHandler(*appModules.user, logger)

	r.Get("/health", healthHandler) // GET /health
	r.Route("/api/v1", func(r chi.Router) {
		// URLs unauthenticated
		r.Get("/health", healthHandler) // GET /api/v1/health
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", authHandler.Login) // POST /api/v1/auth/signin
		})
		// swagger: http://localhost:5080/api/v1/doc/index.html
		r.Get("/doc/doc.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "src/docs/swagger.json")
		})
		r.Get("/doc/*", httpSwagger.Handler(
			httpSwagger.URL("doc.json"),
		))

		// URLs authenticated via jwt bearer token
		r.With(netw.JwtMiddleware(logger)).Route("/user", func(r chi.Router) {
			r.Get("/{userId}", userHandler.GetUserById) // GET /api/v1/user/u/{userId}
			r.Get("/", userHandler.GetUsers)            // GET /api/v1/user
		})
	})
}
