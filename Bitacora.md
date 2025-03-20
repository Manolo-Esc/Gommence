

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
- docker with postgres
- ejecutar todos los tests

## Featuring
- Hexagonal arch
- data validation
- database integration tests
- e2e tests
- Swagger

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