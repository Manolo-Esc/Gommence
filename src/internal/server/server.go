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

func WebServiceFactory(
	appModules *AppModules,
	logger logger.LoggerService,
	db *gorm.DB,
) http.Handler {
	r := chi.NewRouter()
	// Middlewares globales
	//r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// ver los ejemplos y el codigo de https://github.com/riandyrn/otelchi/metric para capturar metricas sobre las peticiones
	//r.Use(otelchi.Middleware("opomatic", otelchi.WithChiRoutes(r)))  // xxxx  (si lo activamos posiblemente no necesitemos el middleware de log)
	r.Use(netw.LogMiddleware(logger))
	r.Use(netw.NoCacheMiddleware)

	addRoutes(
		appModules,
		r,
		logger,
		db,
	)
	var handler http.Handler = r
	return handler
}

type Config struct {
	Host string
	Port string
}

func initTracerProvider() (*trace.TracerProvider, error) {
	// Exportador OTLP binario (puede configurarse para Jaeger, Prometheus, etc.)
	// ctx := context.Background()
	// exporter, err := otlptrace.New(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	// Crear el proveedor de trazas
	tp := trace.NewTracerProvider(
		//trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(0.2))), XXX configurar esto? Por defecto usa AlwaysSample, guarda el 100% de las trazas
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opomatic"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

func Run(ctx context.Context, args []string, getenv func(string) string, stdin io.Reader, stdout, stderr io.Writer) error {
	// crear un contexto que se cancela con SIGINT, SIGTERM o SIGHUP
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP) // No debemos capturar SIGKILL ni SIGSTOP
	defer cancel()                                                                          // funcion de cancelacion del contexto

	tp, err := initTracerProvider()
	if err != nil {
		fmt.Println("Error initializing OpenTelemetry:", err)
		return err
	}

	// Inicialización de los recursos
	logger := logger.GetLogger()
	defer logger.Sync()

	//docker run --name Test -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_DB=mydb -e POSTGRES_PASSWORD=secret -d postgres:16.3
	dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	//var db *gorm.DB = nil
	logger.Info("Conexión a la base de datos exitosa")
	err = database.Migrate(ctx, db)
	if err != nil {
		log.Fatal("Error migrating or cheking database version: ", err)
	}

	appModules := ProductionAppModulesFactory(logger, db, cache.GetCache())

	// Migraciones automáticas (solo para desarrollo)
	//db.AutoMigrate(&User{})

	config := Config{Host: "127.0.0.1", Port: "5080"}
	// tenantsStore := NewTenantsStore()
	// slackLinkStore := NewSlackLinkStore()
	// msteamsLinkStore := NewMSTeamsLinkStore()
	// proxy := NewProxy()

	// Crear el servidor HTTP
	srv := WebServiceFactory(
		appModules,
		logger,
		db,
	// config,
	// tenantsStore,
	// slackLinkStore,
	// msteamsLinkStore,
	// proxy,
	)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	// Lanzar el servidor en una goroutine para no bloquear
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Manejo de cierre y limpieza
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()                                                                     // Esperar señal de cierre
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // crear un nuevo contexto con un tiempo límite
		defer cancel()
		log.Println("shutting down http server")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
		shutdownCtx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second) // crear un nuevo contexto con un tiempo límite
		defer cancel2()
		log.Println("shutting down OpenTelemetry")
		if err := tp.Shutdown(shutdownCtx2); err != nil {
			fmt.Fprintf(stderr, "error shutting down OpenTelemetry: %s\n", err)
		}
	}()
	// Esperar que todas las goroutines terminen
	wg.Wait()
	return nil
}
