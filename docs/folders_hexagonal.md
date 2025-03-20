



# Estructura de Carpetas Recomendada para arquitectura hexagonal
La idea es dividir en **dominio, puertos y adaptadores**, manteniendo una separación clara entre la lógica de negocio y la infraestructura.

```
/project-root
│── /cmd               # Punto de entrada de la app (main.go)
│── /internal          # Código interno del proyecto (no expuesto)
│   │── /domain        # 📌 Entidades y lógica de negocio, modelos de dominio puros (no ORM)
│   │   │── user.go    # Entidad User
│   │   │── invoice.go # Entidad Invoice
│   │
│   │── /dto            # 📌 DTOs (estructuras para entrada/salida de datos)
│   │   │── user.go    # Entidad User
│   │
│   │── /app           # 📌 Casos de uso / Servicios de negocio
│   │   │── user_service.go    # Casos de uso de User
│   │   │── invoice_service.go # Casos de uso de Invoice
│   │
│   │── /ports         # 📌 Interfaces (puertos)
│   │   │── user.go    # UserRepository, UserService interfaces
│   │   │── invoice.go # InvoiceRepository, InvoiceService interfaces
│   │
│   │── /adapters      # 📌 Implementaciones de interfaces
│   │   │── /db
│   │   │   │── user_repository.go      # Implementa UserRepository
│   │   │   │── user_model.go           # Entidad ORM User
│   │   │   │── invoice_repository.go   # Implementa InvoiceRepository
│   │   │── /http
│   │   │   │── user_handler.go       # Implementa HTTP handlers
│   │   │   │── invoice_handler.go    # Implementa HTTP handlers
│   │
│   │── /config        # 📌 Configuración de la aplicación
│   │── /infra         # 📌 Código de infraestructura (db, cache, etc.)
│   │── /tests         # 📌 Pruebas unitarias e integración
│
│── /pkg               # Código reutilizable si es necesario
│── /vendor            # Dependencias externas (go mod vendor)
│── go.mod             # Módulo de Go
│── go.sum             # Checksum de dependencias
```

---


- **Interfaces (Puertos)** → En `/internal/ports/`
  - Define las **interfaces** que el dominio necesita para interactuar con la infraestructura.
  - Ejemplo (`ports/user.go`):
    ```go
    package ports

    import "project-root/internal/domain"

    type UserRepository interface {
        GetByID(id string) (*domain.User, error)
        Create(user *domain.User) error
    }

    type UserService interface {
        RegisterUser(name string, email string) (*domain.User, error)
    }
    ```

- **Adaptadores (Implementaciones de los puertos)** → En `/internal/adapters/`
  - Cada implementación (DB, HTTP, gRPC, etc.) va en su propio subpaquete.
  - Ejemplo (`adapters/db/user_repository.go`):
    ```go
    package db

    import (
        "project-root/internal/domain"
        "project-root/internal/ports"
    )

    type UserRepositoryDB struct {
        db *sql.DB
    }

    func NewUserRepository(db *sql.DB) ports.UserRepository {
        return &UserRepositoryDB{db: db}
    }

    func (r *UserRepositoryDB) GetByID(id string) (*domain.User, error) {
        // Lógica de acceso a la base de datos
    }
    ```

- **Ejemplo de DTO (`internal/dto/user.go`)**
```go
package dto

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

type CreateUserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

- **Ejemplo de entidad en `/internal/domain/user.go`**  
```go
package domain

import "github.com/google/uuid"

type User struct {
    ID    string
    Name  string
    Email string
}

func NewUser(name, email string) *User {
    return &User{
        ID:    uuid.New().String(),
        Name:  name,
        Email: email,
    }
}
```


- **Ejemplo de modelo GORM (`internal/adapters/db/user_model.go`)**
```go
package db

import "gorm.io/gorm"

type UserModel struct {
    gorm.Model
    Name  string `gorm:"size:255"`
    Email string `gorm:"uniqueIndex"`
}
```
Y en el repositorio, conviertes entre la entidad de dominio y el modelo del ORM, por ejemplo:
```go
func (r *UserRepositoryDB) GetByID(id string) (*domain.User, error) {
    var userModel UserModel
    if err := r.db.First(&userModel, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &domain.User{
        ID:    userModel.ID,
        Name:  userModel.Name,
        Email: userModel.Email,
    }, nil
}
```
