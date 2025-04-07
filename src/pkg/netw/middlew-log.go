package netw

import (
	"fmt"
	"net/http"

	"github.com/Manolo-Esc/gommence/src/pkg/logger"
)

func LogMiddleware(logger logger.LoggerService) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			nextHandler.ServeHTTP(w, r)
		})
	}
}
