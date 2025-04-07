package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/infra/database"
	"github.com/Manolo-Esc/gommence/src/pkg/cache"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/Manolo-Esc/gommence/src/pkg/netw"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func teste(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	fmt.Fprintf(w, "userId: %s", id)
}

func WebServiceFactory(appModules *AppModules, logger logger.LoggerService, db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	// Global Middlewares
	//r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// See samples in https://github.com/riandyrn/otelchi/metric to record metrics about the received calls
	r.Use(netw.LogMiddleware(logger))
	r.Use(netw.NoCacheMiddleware)
	addRoutes(appModules, r, logger, db)
	var handler http.Handler = r
	return handler
}

type Config struct {
	Host string
	Port string
}

func initTracerProvider() (*trace.TracerProvider, error) {
	// Binary OTLP exporter (can be configured for Jaeger, Prometheus, etc.)
	// ctx := context.Background()
	// exporter, err := otlptrace.New(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// Console exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		//trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.2))), // Default is AlwaysSample, keeping 100% traces
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("gommence"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func initDatabase(ctx context.Context, dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return nil, err
	}
	err = database.Migrate(ctx, db)
	if err != nil {
		log.Fatal("Error migrating or cheking database version: ", err)
		return nil, err
	}
	return db, nil
}

func Run(ctx context.Context, args []string, getenv func(string) string, stdin io.Reader, stdout, stderr io.Writer) error {
	// Create a context that can be cancelled with SIGINT, SIGTERM o SIGHUP
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP) // We must not capture SIGKILL or SIGSTOP
	defer cancel()

	tp, err := initTracerProvider()
	if err != nil {
		fmt.Println("Error initializing OpenTelemetry:", err)
		return err
	}
	logger := logger.GetLogger()
	defer logger.Sync()

	db, err := initDatabase(ctx, "host=localhost user=postgres password=secret dbname=my_db port=5432 sslmode=disable TimeZone=UTC") // args or getenv should be used here
	if err != nil {
		return err
	}

	appModules := ProductionAppModulesFactory(logger, db, cache.GetCache())

	config := Config{Host: "127.0.0.1", Port: "5080"} // args or getenv should be used here

	srv := WebServiceFactory(appModules, logger, db)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	// launch server in a goroutine to avoid blocking this one
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { // cleaning goroutine
		defer wg.Done()
		<-ctx.Done() // wait closing signal
		// Program shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // new context with timeout
		defer cancel()
		log.Println("shutting down http server")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
		shutdownCtx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second) // new context with timeout
		defer cancel2()
		log.Println("shutting down OpenTelemetry")
		if err := tp.Shutdown(shutdownCtx2); err != nil {
			fmt.Fprintf(stderr, "error shutting down OpenTelemetry: %s\n", err)
		}
	}()
	// Wait for cleaning goroutine
	wg.Wait()
	return nil
}
