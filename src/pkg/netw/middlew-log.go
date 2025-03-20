package netw

import (
	"fmt"
	"net/http"

	"github.com/Manolo-Esc/gommence/src/pkg/logger"
)

/* version net.http
func LogMiddleware(nextHandler http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if !currentUser(r).IsAdmin {
		// 	http.NotFound(w, r)
		// 	return
		// }
		logger.Info("Llamada recibida")
		nextHandler.ServeHTTP(w, r)
	})
}
*/

// version chi
func LogMiddleware(logger logger.LoggerService) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			nextHandler.ServeHTTP(w, r)
		})
	}
}
