package repos_db

import (
	"github.com/Manolo-Esc/gommence/src/internal/ports"
)

type PermissionRepositoryDB struct {
	dbInfra *DBReposInfra
}

func NewPermissionRepository(dbInfra *DBReposInfra) ports.PermissionRepository {
	return &PermissionRepositoryDB{dbInfra: dbInfra}
}

// func (r *PermissionRepositoryDB) GetByID(ctx context.Context, id string) (*domain.User, error) {

// 	return nil, nil
// }
