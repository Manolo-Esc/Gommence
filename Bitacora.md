

## To Do
- go mod tidy
- verificar que no hay ningun folder logs que vaya a git
- limpiar folder docs
- cmd/ws_test.go
- alguien usa 
  - domain/permission.go?? If so, ver cual de las dos implementaciones y borrar la otra
  - mocks/permission_mocks
- opo_uid
- libtest.go: limpiar
- crear swagger
  - ver si en router.go hace falta poner los tag.name
- docker with postgres: //docker run --name Test -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_DB=my_db -e POSTGRES_PASSWORD=secret -d postgres:16.3
- ejecutar todos los tests
- tool para crear typescripts?

## Featuring
- Hexagonal arch
- data validation
- database integration tests
- e2e tests
- Swagger (http://localhost:5080/api/v1/doc/index.html)

### Misc
- documentar que todos pueden devolver http.StatusUnauthorized??


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






## Procedimientos Back
- ejecutar back: en root ejecutar: `go run src/cmd/main.go` o en src: `go run cmd/main.go`
- descargar dependencias: `go mod download`
- generar swagger: cambiar a src y ejecutar `swag init -g router.go -d cmd,internal/dtos,internal/adapters/rest`
- Tests
  - ejecutar tests: en root o src ejecutar: `go test ./...` o en modo verbose `\ `
  - ejecutar tests de un paquete en concreto: ir al paquete y ejecutar `go test` o `go test <ruta_al_paquete>`
  - informe de cobertura de tests: en root o src ejecutar: go test -cover ./...
  - detector de race conditions: en root o src ejecutar: go test -race ./...
- Mocks
  - Cambiar al folder ports (o donde esté el interface del que se quiera generar mocks, pero atención al nombre en destination no colisione)
  - `mockgen -source=entity_ports.go -destination=../mocks/entity_mocks.go -package=mocks`



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




### Comparar los tsconfig.json de los 3 proyectos? 

Si quieres configurar un tsconfig.json global para todo el monorepo, puedes crear uno en la raíz del proyecto para que los otros workspaces puedan extenderlo. Este archivo será común para todos los workspaces, y cada uno puede sobrescribirlo si es necesario.

json

{
  "compilerOptions": {
    "baseUrl": ".",  // Permite resolver los paths relativos desde la raíz
    "paths": {
      "@my-monorepo/*": ["packages/*/src"]
    },
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true
  },
  "exclude": ["node_modules", "dist"]
}


