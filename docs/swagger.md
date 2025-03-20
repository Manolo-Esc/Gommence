
# [`swag`](https://github.com/swaggo/swag). 

Esta herramienta permite agregar **anotaciones en el código** y luego generar la documentación en formato Swagger. La documentación se genera con `swag init` y se sirve con `http-swagger`. 

1. Instalar el comando `swag`  
 ```sh
   go install github.com/swaggo/swag/cmd/swag@latest
   ```


2. Instalar los middlewares y paquetes necesarios
 ```sh
   go get -u github.com/swaggo/http-swagger
   go get -u github.com/swaggo/swag
   ```

3. Agregar anotaciones en el código
   Formato compatible con [swaggo](https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format), por ejemplo:

   ```go
   package main

   import (
       "net/http"

       "github.com/go-chi/chi/v5"
       "github.com/swaggo/http-swagger"
       _ "github.com/tuusuario/tuproject/docs" // Importar el paquete generado
   )

   // @title           API de Ejemplo con Chi
   // @version         1.0
   // @description     Esta es una API de ejemplo documentada con Swagger en Chi.
   // @host           localhost:8080
   // @BasePath       /api

   func main() {
       r := chi.NewRouter()

       r.Get("/swagger/*", httpSwagger.WrapHandler)

       r.Get("/api/saludo", saludoHandler)

       http.ListenAndServe(":8080", r)
   }

   //
   // IMPORTANTE!   IMPORTANTE!    IMPORTANTE!    IMPORTANTE! 
   //
   // NO DEJAR ESPACIOS EN BLANCO ENTRE LAS ANOTACIONES Y LA FUNCION!!!
   //
   // IMPORTANTE!   IMPORTANTE!    IMPORTANTE!    IMPORTANTE! 
   //

   // saludoHandler responde con un mensaje de saludo.
   // @Summary Saludo de bienvenida
   // @Description Responde con un mensaje de bienvenida.
   // @Tags Saludo
   // @Accept  json
   // @Produce  json
   // @Success 200 {object} map[string]string
   // @Router /api/saludo [get]
   func saludoHandler(w http.ResponseWriter, r *http.Request) {
       w.Header().Set("Content-Type", "application/json")
       w.Write([]byte(`{"mensaje": "¡Hola, mundo!"}`))
   }
   ```

4. **Generar la documentación**  
   Ejecuta el siguiente comando en el directorio raíz del proyecto:
```sh
   swag init
   ```
   Esto generará una carpeta **`docs/`** con los archivos `swagger.json` y `swagger.yaml`.

5. **Servir la documentación Swagger**  
   Esto ha sido mas complicado que lo que cuentan aqui. Ver el codigo 
   Agrega el **handler de Swagger UI** en el router:
   ```go
   r.Get("/swagger/*", httpSwagger.WrapHandler)
   ```
   Luego, inicia el servidor y accede a:
   ```
   http://localhost:8080/swagger/index.html
   ```
6.  Customizar el aspecto. Ver: 
   - https://github.com/swaggo/http-swagger
   - https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/
