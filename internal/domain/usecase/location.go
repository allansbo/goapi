package usecase

import (
	"github.com/allansbo/goapi/internal/app/server/dto"
	"github.com/allansbo/goapi/internal/config"
	"github.com/allansbo/goapi/internal/domain/entity"
	"github.com/allansbo/goapi/internal/provider/db"
)

var (
	AppConfig  *config.EnvConfig
	repository db.Repository
	err        error
)

// LoadAppConfig loads the application configuration from environment variables.
func LoadAppConfig() error {
	AppConfig, err = config.LoadEnvConfig()
	if err != nil {
		return err
	}

	return nil
}

// LoadDatabaseRepository initializes the database repository using the application configuration.
func LoadDatabaseRepository() error {
	repository = db.NewMongoDBRepository(AppConfig)
	if err := repository.Ping(); err != nil {
		return err
	}

	return nil
}

// SaveLocation saves a new location in the database and returns the saved location.
// It takes a pointer to dto.LocationInApp as input, which contains the validated location data.
// It returns a pointer to dto.LocationOutApp and an error if any occurs.
func SaveLocation(locationDataIn *dto.LocationInApp) (*dto.LocationOutApp, error) {
	locationEntity := entity.NewLocationInApp(locationDataIn)
	locationOutDB := locationEntity.NewLocationOutDB()

	var err error
	locationEntity.ID, err = repository.InsertOne(locationOutDB)
	if err != nil {
		return nil, err
	}

	return locationEntity.NewLocationOutApp(), nil
}

// GetLocationById retrieves a location by its ID from the database.
// It takes a string ID as input and returns a pointer to dto.LocationOutApp and an error if any occurs.
func GetLocationById(id string) (*dto.LocationOutApp, error) {
	locationInDB, err := repository.GetOne(id)
	if err != nil {
		return nil, err
	}

	locationEntity := entity.NewLocationInDB(locationInDB)

	return locationEntity.NewLocationOutApp(), nil
}

// UpdateLocation updates an existing location in the database.
// It takes a string ID and a pointer to dto.LocationInApp as input,
// which contains the validated location data that will be updated.
// It returns a boolean indicating success and an error if any occurs.
func UpdateLocation(id string, locationDataIn *dto.LocationInApp) (bool, error) {
	locationEntity := entity.NewLocationInApp(locationDataIn)
	locationOutDB := locationEntity.NewLocationOutDB()

	res, err := repository.UpdateOne(id, locationOutDB)
	if err != nil {
		return false, err
	}

	return res, nil
}

// DeleteLocation deletes a location by its ID from the database.
// It takes a string ID as input and returns a boolean indicating success and an error if any occurs.
func DeleteLocation(id string) (bool, error) {
	res, err := repository.DeleteOne(id)
	if err != nil {
		return false, err
	}
	return res, nil
}
