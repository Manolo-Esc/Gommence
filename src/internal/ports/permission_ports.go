package ports

import (
	"context"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
)

type PermissionRepository interface {
	//CreateCourse(ctx context.Context, creationParams *dtos.CourseCreation) (*dtos.Course, error)
}

type PermissionService interface {
	IsSameUserOrHasSomePermission(byUser string, forUser string, permissions []domain.Permission) (bool, APIError)
	GetUserGlobalPermissions(ctx context.Context, forUser string, byUser string) ([]domain.Permission, APIError)
}
