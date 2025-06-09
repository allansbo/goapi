package db

import (
	"github.com/allansbo/goapi/internal/app/server/dto"
)

// Repository defines the interface for database operations related to locations.
type Repository interface {
	Connect() error
	InsertOne(location *dto.LocationOutDB) (string, error)
	GetOne(id string) (*dto.LocationInDB, error)
	UpdateOne(id string, location *dto.LocationOutDB) (bool, error)
	DeleteOne(id string) (bool, error)
}
