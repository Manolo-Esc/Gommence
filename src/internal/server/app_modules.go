package server

import (
	"github.com/Manolo-Esc/gommence/src/internal/adapters/repos_db"
	"github.com/Manolo-Esc/gommence/src/internal/app"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/cache"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"gorm.io/gorm"
)

type AppModules struct {
	auth       *ports.AuthService
	permission *ports.PermissionService
	user       *ports.UserService
}

func ProductionAppModulesFactory(logger logger.LoggerService, db *gorm.DB, cache cache.CacheService) *AppModules {
	dbInfra := repos_db.DBReposInfra{
		Db:     db,
		Logger: logger,
	}
	permission := app.NewPermissionService(repos_db.NewPermissionRepository(&dbInfra), cache, logger)
	serviceInfra := app.ServiceInfra{
		Logger:      logger,
		Cache:       cache,
		Permissions: permission,
	}
	user := app.NewUserService(repos_db.NewUserRepository(&dbInfra), &serviceInfra)
	auth := app.NewAuthService(&serviceInfra, user)
	return &AppModules{
		auth:       &auth,
		permission: &permission,
		user:       &user,
	}
}
