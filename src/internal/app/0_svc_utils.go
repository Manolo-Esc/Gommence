package app

import (
	"github.com/Manolo-Esc/gommence/src/internal/mocks"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/cache"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"go.uber.org/mock/gomock"
)

type ServiceInfra struct {
	Logger      logger.LoggerService
	Cache       cache.CacheService
	Permissions ports.PermissionService
}

func mockServiceInfra(ctrl *gomock.Controller) *ServiceInfra {
	return &ServiceInfra{
		Logger:      logger.GetNopLogger(),
		Cache:       cache.NewCache(),
		Permissions: mocks.NewMockPermissionService(ctrl),
	}
}
