
### go 

- Ejecuta el archivo principal desde su ubicación con `go run src/cmd/*.go` 
- si deseas compilarlo, usa go build -o my-app cmd/main.go.

- Investigar
  - .gitignore. 
    - Se incluye el .sum?

### Learn
- https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
- https://lets-go.alexedwards.net/sample/02.04-wildcard-route-patterns.html
- https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/hello-world


### go vs npm

**Resumen de equivalencias**
| Acción en Node.js        | Acción en Go              |
|------------------------|-----------------------------|
| `package.json`         | `go.mod`                    |
| `package-lock.json`    | `go.sum`                    |
| `npm init -y`          | `go mod init nombre-modulo` |
| `npm install paquete`  | `go get paquete`            |
| `npm ci`               | `go mod download`           |
| `npm install`          | `go mod tidy`               |
| `npm update`           | `go get -u ./...`           |
| `npm prune`            | `go mod tidy`               |
| `node index.js`        | `go run main.go`            |
| `npm run build`        | `go build`                  |
| `npm unistall`         | dejar de importar el paquete y hacer un `go mod tidy`|

Go es más simple en gestión de dependencias porque no necesita un equivalente exacto a `node_modules`, ya que las dependencias se almacenan en caché globalmente (`$GOPATH/pkg/mod`).


