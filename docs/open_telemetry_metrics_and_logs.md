

# Métricas con OpenTelemetry

OpenTelemetry soporta métricas como **counters**, **histogramas**, **gauge** y **summary**.

#### 📝 Pasos básicos para integrar métricas:

1. **Instalar los paquetes de métricas**:

   ```sh
   go get go.opentelemetry.io/otel/sdk/metrics
   go get go.opentelemetry.io/otel/exporters/otlp/otlpmetrics
   go get go.opentelemetry.io/otel/metric
   ```

2. **Configurar el proveedor de métricas**:

   Primero, debes crear un `MeterProvider` que gestione la exportación de métricas:

   ```go
   package main

   import (
       "context"
       "fmt"

       "go.opentelemetry.io/otel"
       "go.opentelemetry.io/otel/sdk/metrics"
       "go.opentelemetry.io/otel/sdk/resource"
       "go.opentelemetry.io/otel/exporters/otlp/otlpmetrics"
   )

   func initMetricsProvider() (*metrics.MeterProvider, error) {
       // Exportador OTLP para métricas
       exporter, err := otlpmetrics.New(context.Background())
       if err != nil {
           return nil, fmt.Errorf("failed to create OTLP metrics exporter: %w", err)
       }

       // Configurar el MeterProvider
       mp := metrics.NewMeterProvider(
           metrics.WithBatcher(exporter),
           metrics.WithResource(resource.NewWithAttributes(
               "service.name", "mi-api-metrics",
           )),
       )

       // Registrar el MeterProvider
       otel.SetMeterProvider(mp)
       return mp, nil
   }
   ```

3. **Crear e instrumentar métricas**:

   Con el `MeterProvider` configurado, puedes crear métricas de tipo **counter**, **gauge** o **histograma**:

   ```go
   package main

   import (
       "context"
       "fmt"
       "go.opentelemetry.io/otel"
       "go.opentelemetry.io/otel/metric"
       "go.opentelemetry.io/otel/label"
   )

   func recordMetrics() {
       meter := otel.Meter("mi-api-metrics")

       // Crear un contador
       requestsCounter := meter.NewInt64Counter("http_requests_total")

       // Registrar la métrica
       requestsCounter.Add(context.Background(), 1, label.Key("method").String("GET"))

       // Crear un histograma
       responseTimeHistogram := meter.NewFloat64Histogram("http_response_duration_seconds")
       responseTimeHistogram.Record(context.Background(), 0.123, label.Key("method").String("GET"))
   }
   ```

3. **Caso Prometheus**   
Hay que conectar el exporter con el endpoint /metrics:
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

func main() {
	// Crea el exporter para Prometheus
	exporter, err := prometheus.New(prometheus.WithoutDescriptorNames()) // Sin nombres duplicados en los descriptores
	if err != nil {
		log.Fatalf("failed to create Prometheus exporter: %v", err)
	}

	// Configura el MeterProvider
	meterProvider := metric.NewMeterProvider(
		metric.WithExporter(exporter),
	)
	otel.SetMeterProvider(meterProvider)

	// Crea el Meter
	meter := meterProvider.Meter("my-meter")

	// Crea un contador de métricas
	counter, err := meter.SyncInt64().Counter("my_counter")
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	// Actualiza la métrica con un valor
	counter.Add(context.Background(), 1, attribute.String("status", "success"))

	// Exponer las métricas a Prometheus a través de /metrics
	http.HandleFunc("/metrics", exporter.ServeHTTP)
	log.Println("Serving metrics at :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

4. **Visualización de las métricas**:

   Las métricas se exportan a un backend como **Prometheus**, **Jaeger** o **OTLP**. Si usas **Prometheus**, el backend tiene que consultar las métricas de la aplicación (por ejemplo, usando un endpoint `/metrics`).

---

# Logs con OpenTelemetry

Los logs son cruciales para la depuración y el seguimiento de eventos. OpenTelemetry permite capturar logs estructurados y enviarlos a tu sistema de logging preferido.

#### 📝 Pasos básicos para integrar logs:

1. **Instalar paquetes para logging**:

   ```sh
   go get go.opentelemetry.io/otel/sdk/logs
   go get go.opentelemetry.io/otel/exporters/otlp/otlpLogs
   ```

2. **Configurar el `LogProvider`**:

   Similar a las métricas, necesitas configurar un **LogProvider** para gestionar los logs:

   ```go
   package main

   import (
       "context"
       "fmt"
       "go.opentelemetry.io/otel"
       "go.opentelemetry.io/otel/sdk/logs"
       "go.opentelemetry.io/otel/exporters/otlp/otlpLogs"
   )

   func initLogProvider() (*logs.LogProvider, error) {
       // Exportador OTLP para logs
       exporter, err := otlpLogs.New(context.Background())
       if err != nil {
           return nil, fmt.Errorf("failed to create OTLP logs exporter: %w", err)
       }

       // Configurar el LogProvider
       lp := logs.NewLogProvider(
           logs.WithBatcher(exporter),
       )

       // Registrar el LogProvider
       otel.SetLogProvider(lp)
       return lp, nil
   }
   ```

3. **Generar logs estructurados**:

   Para generar logs, puedes usar el **Logger** que OpenTelemetry proporciona, y agregarle información relevante a los eventos de tu aplicación:

   ```go
   package main

   import (
       "context"
       "fmt"
       "go.opentelemetry.io/otel"
       "go.opentelemetry.io/otel/attribute"
   )

   func generateLogs() {
       logger := otel.GetLogger("mi-api-logs")

       // Crear un log estructurado
       logger.Info("Request received", attribute.String("method", "GET"), attribute.String("path", "/api"))

       // Crear un log con nivel de error
       logger.Error("Failed to process request", attribute.String("error", "database timeout"))
   }
   ```

4. **Exportar los logs**:

   Los logs se pueden exportar a un sistema como **Jaeger**, **OTLP**, o cualquier otro backend de logging como **Elasticsearch**, **Datadog**, **NewRelic**, etc.

---

### 🔹 **3️⃣ Integración de Métricas y Logs en Chi y GORM**

#### **Métricas en Chi**: Puedes contar el número de solicitudes HTTP.

```go
r := chi.NewRouter()
r.Use(otelchi.Middleware("mi-api"))
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    requestsCounter.Add(r.Context(), 1, label.Key("method").String("GET"))
    w.Write([]byte("Hello"))
})
```

#### **Logs en GORM**: Puedes registrar logs cuando se realicen consultas a la base de datos.

```go
db.Callback().Create().After("gorm:create").Register("log_create", func(db *gorm.DB) {
    logger := otel.GetLogger("gorm")
    logger.Info("GORM query executed", attribute.String("query", db.Statement.SQL.String()))
})
```

---

### ✅ **Resumen**

1. **Métricas**: Usa OpenTelemetry para registrar métricas (counter, histograma, etc.) sobre tus servicios, consultas y tráfico HTTP.
2. **Logs**: Registra logs estructurados para ayudar en la depuración, los cuales se pueden exportar a plataformas de logging.
3. **Integración con Chi y GORM**: Las métricas y logs se integran fácilmente en Chi (para HTTP) y GORM (para base de datos).

🔹 **¿Te gustaría ver cómo configurar un backend como Prometheus o Jaeger para visualizar métricas y logs?** 🚀