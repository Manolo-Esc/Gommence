
## Doc
https://github.com/go-chi/


## Install
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware



En este caso, construiremos una **API REST para gestión de usuarios** con:  
✅ Rutas para CRUD (`GET`, `POST`, `PUT`, `DELETE`)  
✅ Middleware de autenticación  
✅ Validación de datos de entrada  
✅ Respuestas JSON  

---

### 📌 **Estructura del proyecto**  
```bash
myapp/
│── main.go
│── handlers/
│   ├── user_handlers.go
│── middleware/
│   ├── auth.go
│── models/
│   ├── user.go
```
- **`main.go`** → Punto de entrada, configura el router.  
- **`handlers/user_handlers.go`** → Lógica de los endpoints.  
- **`middleware/auth.go`** → Middleware de autenticación.  
- **`models/user.go`** → Definición del modelo `User`.

---

### 🏗 **Paso 1: Definir el modelo de usuario**  
📌 **Archivo:** `models/user.go`  
```go
package models

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

---

### 🏗 **Paso 2: Implementar los handlers**  
📌 **Archivo:** `handlers/user_handlers.go`  
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/myapp/models"
)

var users = map[string]models.User{}

// Obtener lista de usuarios
func GetUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// Obtener un usuario por ID
func GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    user, exists := users[userID]
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(user)
}

// Crear un nuevo usuario
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    users[user.ID] = user
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Actualizar un usuario
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    if _, exists := users[userID]; !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    var updatedUser models.User
    if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    users[userID] = updatedUser
    json.NewEncoder(w).Encode(updatedUser)
}

// Eliminar un usuario
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    if _, exists := users[userID]; !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    delete(users, userID)
    w.WriteHeader(http.StatusNoContent)
}
```

---

### 🏗 **Paso 3: Middleware de autenticación**  
📌 **Archivo:** `middleware/auth.go`  
```go
package middleware

import (
    "net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token != "Bearer mysecrettoken" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

### 🏗 **Paso 4: Configurar el router en `main.go`**  
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
    r.Use(middleware.RealIP)
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

### 🚀 **Probando la API con `curl`**
#### ➜ **Crear un usuario**
```sh
curl -X POST http://localhost:8080/api/users \
     -H "Authorization: Bearer mysecrettoken" \
     -H "Content-Type: application/json" \
     -d '{"id":"1","name":"Alice","email":"alice@example.com"}'
```
#### ➜ **Listar usuarios**
```sh
curl -X GET http://localhost:8080/api/users \
     -H "Authorization: Bearer mysecrettoken"
```
#### ➜ **Obtener un usuario por ID**
```sh
curl -X GET http://localhost:8080/api/users/1 \
     -H "Authorization: Bearer mysecrettoken"
```
#### ➜ **Actualizar un usuario**
```sh
curl -X PUT http://localhost:8080/api/users/1 \
     -H "Authorization: Bearer mysecrettoken" \
     -H "Content-Type: application/json" \
     -d '{"id":"1","name":"Alice Updated","email":"alice.new@example.com"}'
```
#### ➜ **Eliminar un usuario**
```sh
curl -X DELETE http://localhost:8080/api/users/1 \
     -H "Authorization: Bearer mysecrettoken"
```

---

### ✅ **¿Qué hemos logrado?**
🔹 Rutas RESTful para gestión de usuarios.  
🔹 Middleware de autenticación.  
🔹 Manejo de JSON en requests/responses.  
🔹 Uso de `chi` para organizar rutas y middlewares.  
