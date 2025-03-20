package users

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func HandleUser(logger *zap.Logger) http.Handler {
	//thing := prepareThing()  // closure with initialization data for the handlers. Note that is read only, mutexes are needed if it is modified
	// can also be used to manage custom data types that will not be used anywhere else
	// type request struct {
	// 	Name string
	// }
	// type response struct {
	// 	Greeting string `json:"greeting"`
	// }
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// use thing to handle request
			logger.Info("HandleUser - handleSomething")

			ctx := r.Context() // Extraer el contexto de la solicitud
			// Establecer un timeout de 1 segundo para la operación
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			// Obtener el parámetro "id" desde la URL
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "missing id", http.StatusBadRequest)
				return
			}

			// Llamar a una función "hija" que reciba el contexto
			data, err := fetchData(ctx, id)
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					http.Error(w, "request timed out", http.StatusGatewayTimeout)
				} else {
					http.Error(w, "failed to fetch data", http.StatusInternalServerError)
				}
				return
			}
			response := map[string]string{"id": id, "data": data}
			json.NewEncoder(w).Encode(response)
		})
}

func fetchData(ctx context.Context, id string) (string, error) {
	// Simular una operación de larga duración
	select {
	case <-time.After(2 * time.Second):
		return "data for " + id, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
