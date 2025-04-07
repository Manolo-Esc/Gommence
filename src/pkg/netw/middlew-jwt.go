package netw

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Manolo-Esc/gommence/src/internal/infra/jwt"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
)

type contextKey string // to avoid collision with other context keys
const userInfoKey contextKey = "gommence_tk_userInfo"

/*
How to use downstream:

	if userInfo, ok := netw.JwtTokenClaims(ctx); ok {
		fmt.Println("User:", userInfo["user"])
		fmt.Println("Expiration:", userInfo["exp"])
	}
*/
func JwtGetTokenClaims(ctx context.Context) (map[string]string, bool) {
	userInfo, ok := ctx.Value(userInfoKey).(map[string]string)
	return userInfo, ok
}

func JwtGetUserInToken(ctx context.Context) string {
	claims, ok := JwtGetTokenClaims(ctx)
	if ok && claims != nil {
		return claims["user"]
	}
	return ""
}

func JwtMiddleware(logger logger.LoggerService) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized) // the text is used in tests!
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized) // the text is used in tests!
				return
			}
			token := parts[1]
			tokenPayload, err := jwt.ValidateToken(token)
			if err != nil {
				http.Error(w, fmt.Sprintf("error in token: %s", err.Error()), http.StatusUnauthorized) // the text is used in tests!
				return
			}
			ctx := context.WithValue(r.Context(), userInfoKey, tokenPayload)
			nextHandler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
