package app

import (
	"context"
	"net/http"

	"github.com/Manolo-Esc/gommence/src/internal/domain"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/cache"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
)

type PermissionServiceImpl struct {
	repo  ports.PermissionRepository
	cache cache.CacheService
}

func NewPermissionService(repo ports.PermissionRepository, cache cache.CacheService, logger logger.LoggerService) ports.PermissionService {
	return &PermissionServiceImpl{repo: repo, cache: cache}
}

func (s *PermissionServiceImpl) GetUserGlobalPermissions(ctx context.Context, forUser string, byUser string) ([]domain.Permission, ports.APIError) {
	return nil, nil
}

func (s *PermissionServiceImpl) IsSameUserOrHasSomePermission(byUser string, forUser string, permissions []domain.Permission) (bool, ports.APIError) {
	if byUser == forUser {
		return true, nil
	}
	// XXX read permissions from cache
	// XXX read permissions from db and store in cache
	// XXX check if user has any of the permissions
	return false, ports.NewAPIError(http.StatusForbidden, "The data is not accessible")
}
