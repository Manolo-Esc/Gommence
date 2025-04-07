package repos_db

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/infra/opo_uid"
	"github.com/Manolo-Esc/gommence/src/internal/ports"
	"github.com/Manolo-Esc/gommence/src/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type BaseDBModel struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // Postgres error code for unique constraint violation
	}
	return false
}

// Creates a new BaseDBModel entity assigning a fancy unique ID
func CreateEntityWithPID[T any](ctx context.Context, db *gorm.DB, entity *T) ports.APIError {
	v := reflect.ValueOf(entity).Elem()
	publicId := v.FieldByName("ID")
	if !publicId.IsValid() {
		return ports.NewAPIError(http.StatusInternalServerError, fmt.Sprintf("type %s does not have 'ID' field", reflect.TypeOf(*entity).Name()))
	}
	if !publicId.CanSet() || publicId.Kind() != reflect.String {
		return ports.NewAPIError(http.StatusInternalServerError, fmt.Sprintf("'ID' field (%d) of type %s can not be set", int(publicId.Kind()), reflect.TypeOf(*entity).Name()))
	}
	currentPublicId := publicId.String()
	if currentPublicId == "" { // we try to respect the publicId if we get one, but we'll change it if we get UniqueViolation error
		publicId.SetString(opo_uid.New())
	}
	times := 0
	for {
		result := db.WithContext(ctx).Create(entity)
		if result.Error == nil {
			break
		}
		if times++; times > 3 { // after several times we have not generated a unique id. It is likely something else is causing the error
			return ports.NewAPIError(http.StatusConflict, result.Error.Error())
		}
		if IsUniqueViolation(result.Error) { // any violation of unique constraint, not just pk
			publicId.SetString(opo_uid.New())
		} else {
			return ports.NewAPIError(http.StatusInternalServerError, result.Error.Error())
		}
	}
	return nil
}

type DBReposInfra struct {
	Db     *gorm.DB
	Logger logger.LoggerService
}
