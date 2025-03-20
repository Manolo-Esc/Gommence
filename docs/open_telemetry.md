

## URLs
- https://opentelemetry.io/docs/instrumentation/go/
- https://pkg.go.dev/go.opentelemetry.io/otel
- https://pkg.go.dev/go.opentelemetry.io/otel/exporters/trace
- https://github.com/open-telemetry/opentelemetry-go
- https://github.com/open-telemetry/opentelemetry-specification

- https://github.com/riandyrn/otelchi/  (middleware chi)




## 🔹 1️⃣ Instalando OpenTelemetry en Go  

Instala los paquetes básicos:  

```sh
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace
go get go.opentelemetry.io/otel/exporters/stdout/stdouttrace
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/propagation
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
go get go.opentelemetry.io/contrib/instrumentation/gorm.io/gorm/otelgorm
go get github.com/riandyrn/otelchi/
```

---

## 🔹 2️⃣ Configurando OpenTelemetry  

Primero, inicializamos el **Tracer Provider** para enviar los datos a **Jaeger, OTLP o cualquier backend de OpenTelemetry**.  

```go
package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
)

func initTracerProvider() (*trace.TracerProvider, error) {
	ctx := context.Background()

	// Exportador OTLP (puedes configurarlo para Jaeger, Prometheus, etc.)
	exporter, err := otlptrace.New(ctx)  // exporta en binario
	if err != nil {
		return nil, err
	}

	// Crear el proveedor de trazas
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("mi-api"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil

	// setup alternativo
	exporter, err := stdouttrace.New(otelstdout.WithPrettyPrint())
    if err != nil {
        log.Fatalf("failed to create stdout exporter: %v", err)
    }

    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
    )
    otel.SetTracerProvider(tp)

}
```

---

## Pasando las cabeceras open telemetry en llamadas a otros servicios
```go
client := http.Client{}
req, _ := http.NewRequest("GET", "http://otro-servicio.com/api", nil)
otel.GetTextMapPropagator().Inject(context.Background(), propagation.HeaderCarrier(req.Header))
client.Do(req)
```

---

## 🔹 3️⃣ Instrumentando Chi para trazar peticiones HTTP  

OpenTelemetry tiene integración nativa con **Chi** mediante `otelchi`.
Esto __añade automáticamente trazas__ a cada request HTTP en Chi

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/github.com/go-chi/chi/otelchi"
)

func main() {
	// Inicializar OpenTelemetry
	tp, err := initTracerProvider()
	if err != nil {
		fmt.Println("Error inicializando OpenTelemetry:", err)
		return
	}
	defer tp.Shutdown(context.Background())

	// Crear router con middleware de OpenTelemetry
	r := chi.NewRouter()

	// Agrega trazas a cada request usando magicamente la configuracion hecha en initTracerProvider
	// Si chi inserta un id de traza, se supone que le llega al gorm para relacionar la llamada a base de datos con la llamada rest
	// Si la llamada rest que recibe chi ya incluye un id de traza, se supone que lo respeta
	r.Use(otelchi.Middleware("mi-api")) 

	// Definir rutas
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Hola, OpenTelemetry en Chi!"))
	})

	// Iniciar servidor
	http.ListenAndServe(":8080", r)
}
```

__Recuperacion de la traza/span en el handler:__  

```go
r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Recuperar el span del contexto de la solicitud
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		// Agregar atributos al span
		span.SetAttributes(attribute.String("custom.attribute", "value"))

		// Loguear con zap, incluyendo el traceID y el spanID
		logger.Info("Handling request",
			zap.String("traceID", span.SpanContext().TraceID.String()),
			zap.String("spanID", span.SpanContext().SpanID.String()),
		)

		// Hacer algo más en el handler (responder la solicitud, etc.)
		w.Write([]byte("Hello, World!"))
	})
```	

__Crear spans hijos__
span y span2 son hermanos o span2 es hijo???
En este ejemplo, se crean spans para representar la operación some-function-operation, y se agrega información adicional mediante atributos.
```go
func someFunction(ctx context.Context) {
	// Crear un span manualmente dentro de una función
	tracer := otel.Tracer("my-tracer")
	ctx, span := tracer.Start(ctx, "some-function-operation")
	span.End()

	// Agregar atributos al span
	span.SetAttributes(attribute.String("operation", "some-function"))

	// Lógica de la operación...
	_, span2 := tracer.Start(ctx, "database-query")
	defer span2.End() // se puede usar defer a nuestra conveniencia, se cerrará el span al final de la funcion
}
```


---

## 🔹 4️⃣ Instrumentando GORM para trazar queries a la base de datos  

Si usas **GORM**, puedes añadir trazas a las consultas usando `otelgorm`:  

```go
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	otelgorm "go.opentelemetry.io/contrib/instrumentation/gorm.io/gorm"
)

func main() {
	// Conectar con SQLite
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Error conectando a la base de datos")
	}

	// Agregar middleware de OpenTelemetry a GORM
	if err := db.Use(otelgorm.NewPlugin()); err != nil { // agrega trazas usando magicamente la configuracion hecha en initTracerProvider
		panic("Error añadiendo OpenTelemetry a GORM")
	}
}
```

📌 **Esto añadirá trazas a cada consulta a la base de datos**.

---

## 🔹 5️⃣ (Opcional) Exportar Trazas a Jaeger  

Si usas **Jaeger**, instala el exportador y configúralo en `initTracerProvider()`:

```sh
go get go.opentelemetry.io/otel/exporters/jaeger
```

Y cambia el exportador:

```go
import (
	"go.opentelemetry.io/otel/exporters/jaeger"
)

exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
```

---

