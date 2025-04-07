package libtest

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

// waitForReady calls the specified endpoint until it gets a 200 response or until the context is cancelled or the timeout is reached.
func WaitForReady(ctx context.Context, timeout time.Duration, endpoint string) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			//fmt.Printf("Error making request: %s\n", err.Error())
			time.Sleep(1000 * time.Millisecond) // wait a little while between checks
			continue
		}
		if resp.StatusCode == http.StatusOK {
			//fmt.Println("Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			time.Sleep(1000 * time.Millisecond) // wait a little while between checks
		}
	}
}

type HttpTestHandler struct {
	Path string
	F    http.Handler
}

type HttpTestHandlerFunc struct {
	Path string
	F    func(http.ResponseWriter, *http.Request)
}

func registerHandlers(mux *http.ServeMux, funcs []HttpTestHandlerFunc, handlers []HttpTestHandler) {
	for _, f := range funcs {
		mux.HandleFunc(f.Path, f.F)
	}
	for _, h := range handlers {
		mux.Handle(h.Path, h.F)
	}
}

// Handler for path "/"
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"mensaje": "Up & Running"}`))
}

func RunSimpleServer(t *testing.T, testFunction []func(t *testing.T, baseURL string), handlersFuncs []HttpTestHandlerFunc, handlers []HttpTestHandler) {
	RunSimpleServerEx(t, testFunction, ":5091", handlersFuncs, handlers)
}

// Sets up a http server with the given handlers and runs the test functions
func RunSimpleServerEx(t *testing.T, testFunctions []func(t *testing.T, baseURL string), port string, handlersFuncs []HttpTestHandlerFunc, handlers []HttpTestHandler) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Timeout para el test
	defer cancel()

	baseURL := fmt.Sprintf("http://localhost%s/", port)
	healthURL := fmt.Sprintf("%shealth", baseURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	registerHandlers(mux, handlersFuncs, handlers)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	//go func() {
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// This is also run in graceful shutdown
		}
	}()
	//}()

	if err := WaitForReady(ctx, 5, healthURL); err != nil {
		t.Fatal("Server is not ready")
	}

	for _, f := range testFunctions {
		f(t, baseURL)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//log.Println("shutting down http server")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			//fmt.Fprintf(stderr, "Error shutting down HTTP server: %s\n", err)
		} else if shutdownCtx.Err() == context.DeadlineExceeded {
			//tf(stderr, "Shutdown failed: timeout exceeded\n")
		} else {
			//log.Println("Server shut down gracefully")
		}
	}()
	wg.Wait()
}
