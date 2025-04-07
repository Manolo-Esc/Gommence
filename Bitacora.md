

## To Do
- xxxx
- verificar que no hay ningun folder "logs" que vaya a git
- limpiar folder docs
- docker with postgres: //docker run --name Test -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_DB=my_db -e POSTGRES_PASSWORD=secret -d postgres:16.3
- docker, docker compose con db?
- run tests
  - ejecutar todos los tests
  - run db tests?
- borrar bitacora.md

textos
"Starter kit"
"Project scaffold"
"Go starter template"

# Gommence

A starter kit for go web services

Introducimos este proyecto para poder tener en marcha rápidamente un servidor web featuring:
- Uso del módulo estándar http y del router [chi](https://github.com/go-chi/) (sin magia)
- Arquitectura hexagonal para testeabilidad y desacoplamiento
- Inyección explícita de dependencias (sin magia)
- Incluye soporte de tokens JWT, hashing de passwords y middleware de autenticación
- Tests unitarios, e2e y de integración con base de datos
- Incluye Dockerfile para una imagen del servicio y docker-compose para ejecutar el servicio con una base de datos
- Integra [testify](https://github.com/stretchr/testify) y [gomock](https://github.com/golang/mock) 
- Integra el ORM [gorm](https://gorm.io) (mucha magia XD)
- Validación de datos con la libreria [validator](https://github.com/go-playground/validator/)
- Incorpora swagger con la libraría [swaggo](https://github.com/swaggo/swag)
- Custom ID generation, much fancier and compact than common GUIDs
- Integra logs con la libraría [zap](https://github.com/uber-go/zap)
- Incluye soporte de OpenTelemetry
- Integra soporte de cache con un layer sobre [ristretto](https://github.com/uber-go/zap) fácilmente extensible a, por ejemplo, redis.


## Featuring
- ideas from 
  https://www.joeshaw.org/error-handling-in-go-http-applications/
  https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
- assertions

### Misc
- documentar que todos pueden devolver http.StatusUnauthorized??

- next steps:
  - tool para crear typescripts?
  - posibilidad de generar metricas propias
  - crear usuario


## Add entities recipe
- definicion de servicio y repo en ports
- Entity, Repo y factory(db) en adapters/repos_db
- dtos en dtos
- servicio y factory(repo) en app
- handlers y factory(servicio) en adapters/rest
- añadir servicio a app_modules.go en server
- añadir handlers en router


## Instructions

### run


### curl
- curl -X POST http://localhost:5080/api/v1/auth/signin \
    -H "Content-Type: application/json" \
    -d '{"email": "user@mail.com", "secret": "password"}'

- curl -X GET http://localhost:5080/api/v1/user \
    -H "Authorization: Bearer your_token"

- curl -X GET http://localhost:5080/api/v1/user/a_valid_id \
    -H "Authorization: Bearer your_token"


## Procedimientos Back
- ejecutar back: en root ejecutar: `go run src/cmd/main.go` o en src: `go run cmd/main.go`
- descargar dependencias: `go mod download`
- vaciar caches de compilacion: `go clean -cache`
- generar swagger: 
  - Instalar: go install github.com/swaggo/swag/cmd/swag@latest
  - cambiar a src y ejecutar `swag init -g router.go -d internal/server,internal/dtos,internal/adapters/rest`
- Tests
  - ejecutar tests: en root o src ejecutar: `go test ./...` o en modo verbose `\ `
  - ejecutar tests de un paquete en concreto: ir al paquete y ejecutar `go test` o `go test <ruta_al_paquete>`
  - informe de cobertura de tests: en root o src ejecutar: go test -cover ./...
  - detector de race conditions: en root o src ejecutar: go test -race ./...
- Mocks
  - Cambiar al folder ports (o donde esté el interface del que se quiera generar mocks, pero atención al nombre en destination no colisione)
  - `mockgen -source=name_of_the_entity_ports.go -destination=../mocks/name_of_the_entity_mocks.go -package=mocks`




### Cheatsheet de la DB del docker
- consola postgreSQL:  
    `psql -dmy_db -Upostgres`

    - Borrar todas las tablas
  ```sql 
        DO $$ DECLARE
            tabname text;
        BEGIN
            FOR tabname IN
                SELECT tablename
                FROM pg_tables
                WHERE schemaname = 'public'
            LOOP
                EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(tabname) || ' CASCADE';
            END LOOP;
        END $$;
``` 
    - Comandos
        \list           lista las databases existentes
        \c savimboplatform   Hacer de savimbo la DB current (para otros comandos)
        \dt             Lista las tablas de la databse current
        \dn             Lista los esquemas de una base de datos
        \d table        Muestra el esquema de la tabla



