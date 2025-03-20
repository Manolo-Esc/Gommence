


go get github.com/go-playground/validator/v10


--

## ✅ **Validación de datos con `go-playground/validator`**  
La librería [`validator`](https://github.com/go-playground/validator) permite agregar reglas de validación a los structs de nuestra API de manera sencilla.

📌 **Modificamos `models/user.go` para incluir validaciones**  
```go
package models

import "github.com/go-playground/validator/v10"

type User struct {
    ID    string `json:"id" validate:"required,uuid4"`  // ID debe ser un UUID válido
    Name  string `json:"name" validate:"required,min=3,max=50"`
    Email string `json:"email" validate:"required,email"`
}

// Validador global (diseñado para ser singleton y thread safe. Si no se usa así se pierde performance. Las funciones no thread-safe están
// marcadas especificamente en la documentacion)
var validate = validator.New(validator.WithRequiredStructEnabled())

// Función para validar un usuario
func (u *User) Validate() error {
    return validate.Struct(u)
}
```
🔹 `validate:"required,uuid4"` → El ID debe ser obligatorio y un UUID válido.  
🔹 `validate:"min=3,max=50"` → Nombre entre 3 y 50 caracteres.  
🔹 `validate:"email"` → Formato de email válido.  

---

## 🚀 **Mejorando el manejo de errores**
En lugar de enviar solo `http.Error()`, creamos una estructura de respuesta JSON para errores.

📌 **Creamos `utils/errors.go`**  
```go
package utils

import (
    "encoding/json"
    "net/http"
)

// Respuesta JSON para errores
type ErrorResponse struct {
    Message string `json:"message"`
}

// Función para enviar errores en formato JSON
// Aqui tenemos dos aproximaciones
// http.Error además "deletes the Content-Length header". Igual habria que incorporarlo si acabamos usando algo de esto
func JSONError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func JSONError(w ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
```

---

## 🔥 **Implementamos validación y manejo de errores en los handlers**  
📌 **Modificamos `handlers/user_handlers.go`**  
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/myapp/models"
    "github.com/myapp/utils"
)

// Obtener usuario por ID
func GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    user, exists := users[userID]
    if !exists {
        utils.JSONError(w, "User not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(user)
}

// Crear usuario con validación
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        utils.JSONError(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // Validar usuario
    if err := user.Validate(); err != nil {
        utils.JSONError(w, err.Error(), http.StatusBadRequest)
        return
    }

    users[user.ID] = user
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Actualizar usuario con validación
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    if _, exists := users[userID]; !exists {
        utils.JSONError(w, "User not found", http.StatusNotFound)
        return
    }

    var updatedUser models.User
    if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
        utils.JSONError(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // Validar datos del usuario
    if err := updatedUser.Validate(); err != nil {
        utils.JSONError(w, err.Error(), http.StatusBadRequest)
        return
    }

    users[userID] = updatedUser
    json.NewEncoder(w).Encode(updatedUser)
}
```
---
## 📌 **Actualizamos `main.go`**
📌 **Archivo:** `main.go`  
```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/myapp/handlers"
    "github.com/myapp/middleware"
)

func main() {
    r := chi.NewRouter()

    // Middlewares globales
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Grupo de rutas protegidas con autenticación
    r.Route("/api", func(r chi.Router) {
        r.With(middleware.AuthMiddleware).Route("/users", func(r chi.Router) {
            r.Get("/", handlers.GetUsers)        // GET /api/users
            r.Post("/", handlers.CreateUser)     // POST /api/users
            r.Get("/{id}", handlers.GetUser)     // GET /api/users/{id}
            r.Put("/{id}", handlers.UpdateUser)  // PUT /api/users/{id}
            r.Delete("/{id}", handlers.DeleteUser) // DELETE /api/users/{id}
        })
    })

    // Iniciar servidor
    http.ListenAndServe(":8080", r)
}
```
---
## 🔥 **Pruebas mejoradas con `curl`**  
🔹 **Intentamos crear un usuario con email inválido**  
```sh
curl -X POST http://localhost:8080/api/users \
     -H "Authorization: Bearer mysecrettoken" \
     -H "Content-Type: application/json" \
     -d '{"id":"123","name":"Al","email":"invalid-email"}'
```
🛑 **Respuesta esperada (400 Bad Request)**  
```json
{
    "message": "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```
🔹 **Intentamos crear un usuario válido**  
```sh
curl -X POST http://localhost:8080/api/users \
     -H "Authorization: Bearer mysecrettoken" \
     -H "Content-Type: application/json" \
     -d '{"id":"550e8400-e29b-41d4-a716-446655440000","name":"Alice","email":"alice@example.com"}'
```
✅ **Respuesta esperada (201 Created)**  
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Alice",
    "email": "alice@example.com"
}
```
🔹 **Intentamos acceder sin token de autenticación**  
```sh
curl -X GET http://localhost:8080/api/users
```
🛑 **Respuesta esperada (401 Unauthorized)**  
```json
{
    "message": "Unauthorized"
}
```
---
## 🏁 **Resumen y beneficios**
✅ **Validación automática con `validator`**  
✅ **Manejo de errores JSON en un solo lugar (`utils/errors.go`)**  
✅ **Middleware de autenticación reutilizable**  
✅ **`chi` permite un router limpio y modular**  

---

### 🎯 **Siguientes pasos**
1️⃣ **Añadir paginación a `GetUsers`**.  
2️⃣ **Usar una base de datos (SQLite, PostgreSQL, etc.) en lugar del mapa `users`**.  
3️⃣ **Usar `envconfig` para manejar configuraciones externas (como el token secreto)**.  

🚀 ¿Te gustaría ver cómo conectar esto con una base de datos o seguir optimizando la API? 😃