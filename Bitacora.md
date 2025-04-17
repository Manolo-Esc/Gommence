

## To Do
- verificar que no hay ningun folder "logs" que vaya a git

- dejamos el middleware-jwt donde está o lo movemos a infra?? no tiene sentido que use infra/jwt, que es menos publico :(
- run tests
  - ejecutar todos los tests
  - run db tests?
- convertir bitacora.md en readme.md

textos
"Starter kit"
"Project scaffold"
"Go starter template"

# Gommence

**A Starter Kit for Go Web Services**

We're introducing this project to quickly spin up a web server featuring:

- Use of the standard `http` module and the [chi](https://github.com/go-chi/) router (no magic included)
- Hexagonal architecture for better decoupling and testability
- Explicit dependency injection (no magic included)
- Built-in support for JWT tokens, password hashing, and authentication middleware
- Unit testing (integrates [testify](https://github.com/stretchr/testify) and [gomock](https://github.com/golang/mock)), plus e2e and database integration tests
- Includes a service Dockerfile and a `docker-compose` setup to run the service with a database
- Integrates the [gorm](https://gorm.io) ORM (loads of magic here, you’ve been warned XD)
- Data validation with the [validator](https://github.com/go-playground/validator/) library
- Swagger integration via [swaggo](https://github.com/swaggo/swag)
- Custom ID generation — much fancier and more compact than your average GUID
- Logging layer on top of [zap](https://github.com/uber-go/zap)
- OpenTelemetry support baked in
- Integrated cache layer using [ristretto](https://github.com/hypermodeinc/ristretto), easily extensible to Redis or others
- Uses [godotenv](https://github.com/joho/godotenv) for environment variable management


## Run with Docker Compose

Assuming you already have Docker installed, run the following from the root folder:

```sh
docker compose up --build -d
```

You can verify the environment is up and running by executing:

```sh
docker ps
```

You should see two containers named *go_web_server* and *go_postgres*. If not, something went sideways.

To check what went wrong, you can inspect the logs:

```sh
docker logs go_web_server
docker logs go_postgres
```

To stop and clean everything up, use:

```sh
docker compose down
```

## Run from Source

You'll need a running Postgres installation with an empty database named *my_db* —or whatever name you’ve set in the `.env` file.

Once you’ve cloned the project, install the dependencies from the project root with:

```sh
go mod download
```

Then run the program from the root directory:

```sh
go run src/cmd/main.go
```

## Calling the Service

Here are a few sample calls you can make to test the service. These examples use _curl_:

### Check Liveness

The _health_ endpoint is available via two routes:

```sh
curl -X GET http://localhost:5080/health
curl -X GET http://localhost:5080/api/v1/health
```

### Get Authentication Token

You can log in and get an authentication token for future requests using the built-in user _user@mail.com_ (with password _password_):

```sh
curl -X POST http://localhost:5080/api/v1/auth/signin \
    -H "Content-Type: application/json" \
    -d '{"email": "user@mail.com", "secret": "password"}'
```

### Get the List of Users

Fetch all users in the system. This request will fail if you don’t provide a valid token (obtained from the previous call):

```sh
curl -X GET http://localhost:5080/api/v1/user \
    -H "Authorization: Bearer the_token_here"
```

### Get Information About a Specific User

To get details about a specific user (whose ID you can retrieve from the previous call), make this request. As before, it requires a valid token:

```sh
curl -X GET http://localhost:5080/api/v1/user/a_valid_id \
    -H "Authorization: Bearer the_token_here"
```

## Hexagonal Architecture

The [hexagonal architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) —also known as the ports and adapters architecture— was proposed by Alistair Cockburn back in 2005. The core idea is to isolate business logic from external concerns like UIs, APIs, storage, or third-party services.

You define **interfaces (ports)** to connect your application’s core logic with its external clients and providers. Then you create **adapters** that either *implement* these interfaces to provide functionality to the core logic (e.g., database access), or *consume* the core logic (e.g., REST endpoints or UIs). Conceptually, you’re splitting your system into **application components** (the business logic) and the **adapters** that talk to each other through ports.

In **Gommence**, business logic lives in the `_app_` and `_domain_` folders. Adapters can be found under `_adapters/repos_db_` and `_adapters/rest_`. The `_dtos_` folder contains the _data transfer objects_ used by the input ports. Interfaces (ports) themselves are defined in the `_ports_` folder.

This level of modularity can seem a bit cumbersome at first. Adding a new entry point to the service might mean creating two new ports, two new adapters, remake the mock objects, and writing the actual logic. Yes, it’s some overhead. But the **decoupling and testability** you get in return? Worth it, especially for large or long-lived projects. If you’re building something that’s going to grow or stick around, this architecture pays off.

### Add entities recipe
**Gommence** incluye entidades _user_ como referencia. Supongamos que ahora queremos añadir una nueva entidad, como *user_post*, or _coche_, o _mascota_. Estos son los puntos que posiblemente haya que tocar:
- definicion de servicio y repo en _ports_
- Creación de entity (representación en base de datos) y repository (database management) en *adapters/repos_db*
- servicio (logica de negocio)) en _app_
- http handlers en _adapters/rest_
- dtos en _dtos_
- añadir servicio a app_modules.go en _server_
- añadir handlers en _router_

## Tests

### Running tests
- Tests
  - ejecutar tests: en root o src ejecutar: `go test ./...` o en modo verbose `\ `
  - ejecutar tests de un paquete (folder) en concreto: ir al paquete y ejecutar `go test` o `go test <ruta_al_paquete>`

### Generate mocks
- Cambiar al folder ports (o donde esté el interface del que se quiera generar mocks, pero atención al nombre en destination no colisione)
- `mockgen -source=name_of_the_entity_ports.go -destination=../mocks/name_of_the_entity_mocks.go -package=mocks`
### Unit tests
### E2E
### Database integration test

## Swagger documentation
The project uses _swaggo_ to generate API documentation from annotations in the code. Puedes ver el resultado apuntando un browser a 
`http://localhost:5080/api/v1/doc/index.html`
Para generar la documentación 
- Instalar _swag_: 
  ```sh
  go install github.com/swaggo/swag/cmd/swag@latest
  ```
- Cambiar al folder src y ejecuta lo que sigue. Debes incluir todos los folders en los que hayas hecho anotaciones en los ficheros go 
  ```sh
  swag init -g router.go -d internal/server,internal/dtos,internal/adapters/rest
  ```

## Créditos
Hasta llegar a este estado este código ha ido tomando con el tiempo ideas de varias fuentes, blogs e IAs, pero quiero destacar especialmente el trabajo de [Mat Ryer](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/) y de [Joe Shaw](https://www.joeshaw.org/error-handling-in-go-http-applications/)
  

### Next steps
  - tool para crear typescripts?
  - posibilidad de generar metricas propias
  - crear usuario


## Cheatsheet de la DB del docker
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
        \dt             Lista las tablas de la current database
        \d table        Muestra el esquema de la tabla
        \list           lista las databases existentes
        \c my_db   Hacer de my_db la current database para otros comandos
       



