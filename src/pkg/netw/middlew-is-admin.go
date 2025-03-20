package netw

import "net/http"

func AdminOnlyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if !currentUser(r).IsAdmin {
		// 	http.NotFound(w, r)
		// 	return
		// }
		h.ServeHTTP(w, r)
	})
}
