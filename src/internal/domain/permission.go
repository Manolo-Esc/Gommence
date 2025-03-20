package domain

import "fmt"

type Permission int16

const (
	PermissionRead   Permission = 1
	PermissionWrite  Permission = 100
	PermissionDelete Permission = 200
	PermissionAdmin  Permission = 300
)

// Definimos un struct con los permisos
type permissions struct {
	Read   Permission
	Write  Permission
	Delete Permission
	Admin  Permission
}

// Constante Permissions con los valores asignados
var Permissions = permissions{
	Read:   1,
	Write:  100,
	Delete: 200,
	Admin:  300,
}

// Map de nombres para el stringer
var permissionNames = map[Permission]string{
	Permissions.Read:   "Read",
	Permissions.Write:  "Write",
	Permissions.Delete: "Delete",
	Permissions.Admin:  "Admin",
}

// Método String para convertir permisos a texto
func (p Permission) String() string {
	if name, ok := permissionNames[p]; ok {
		return name
	}
	return "Unknown"
}

var Roles1 = map[string][]Permission{
	"User":   {Permissions.Read},
	"Editor": {Permissions.Read, Permissions.Write},
	"Admin":  {Permissions.Read, Permissions.Write, Permissions.Delete, Permissions.Admin},
}

func HasPermission1(role string, permission Permission) bool {
	perms, exists := Roles1[role]
	if !exists {
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

func ejemploUso1() {
	fmt.Println("¿Editor puede escribir?", HasPermission1("Editor", Permissions.Write)) // true
	fmt.Println("¿Usuario puede borrar?", HasPermission1("User", Permissions.Delete))   // false
	fmt.Println("¿Admin puede hacer todo?", HasPermission1("Admin", Permissions.Admin)) // true
}

type roles struct {
	User   []Permission
	Editor []Permission
	Admin  []Permission
}

// Constante Roles con los permisos asignados
var Roles2 = roles{
	User:   []Permission{Permissions.Read},
	Editor: []Permission{Permissions.Read, Permissions.Write},
	Admin:  []Permission{Permissions.Read, Permissions.Write, Permissions.Delete, Permissions.Admin},
}

// Función para verificar si un rol tiene un permiso específico
func HasPermission2(role []Permission, permission Permission) bool {
	for _, p := range role {
		if p == permission {
			return true
		}
	}
	return false
}

func ejemploUso2() {
	fmt.Println("¿Editor puede escribir?", HasPermission2(Roles2.Editor, Permissions.Write)) // true
	fmt.Println("¿Usuario puede borrar?", HasPermission2(Roles2.User, Permissions.Delete))   // false
	fmt.Println("¿Admin puede hacer todo?", HasPermission2(Roles2.Admin, Permissions.Admin)) // true
}
