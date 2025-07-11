package db

import (
	"github.com/allansbo/goapi/internal/app/server/dto"
)

// Repository defines the interface for database operations related to locations.
type Repository interface {
	Ping() error
	Stop()
	InsertOne(location *dto.LocationOutDB) (string, error)
	GetOne(id string) (*dto.LocationInDB, error)
	GetAll(query *dto.QueryLocationOutDB) (*dto.QueryLocationInDB, error)
	UpdateOne(id string, location *dto.LocationOutDB) (bool, error)
	DeleteOne(id string) (bool, error)
}
