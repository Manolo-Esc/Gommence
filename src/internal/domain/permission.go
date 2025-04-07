package domain

import "fmt"

type Permission int16

const (
	PermissionRead   Permission = 1
	PermissionWrite  Permission = 100
	PermissionDelete Permission = 200
	PermissionAdmin  Permission = 300
)

type roles struct {
	User   []Permission
	Editor []Permission
	Admin  []Permission
}

var Roles = roles{
	User:   []Permission{PermissionRead},
	Editor: []Permission{PermissionRead, PermissionWrite},
	Admin:  []Permission{PermissionRead, PermissionWrite, PermissionDelete, PermissionAdmin},
}

func HasPermission(role []Permission, permission Permission) bool {
	for _, p := range role {
		if p == permission {
			return true
		}
	}
	return false
}

func usageSample() {
	fmt.Println("Can Editor write?", HasPermission(Roles.Editor, PermissionWrite))     // true
	fmt.Println("Can User delete?", HasPermission(Roles.User, PermissionDelete))       // false
	fmt.Println("Can Admin do anything?", HasPermission(Roles.Admin, PermissionAdmin)) // true
}
