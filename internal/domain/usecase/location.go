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
)

func LoadAppConfig() error {
	var err error
	AppConfig, err = config.LoadEnvConfig()
	if err != nil {
		return err
	}

	return nil
}

func LoadDatabaseRepository() error {
	repository = db.NewMongoDBRepository(AppConfig)
	if err := repository.Connect(); err != nil {
		return err
	}

	return nil
}

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

func GetLocationById(id string) (*dto.LocationOutApp, error) {
	locationInDB, err := repository.GetOne(id)
	if err != nil {
		return nil, err
	}

	locationEntity := entity.NewLocationInDB(locationInDB)

	return locationEntity.NewLocationOutApp(), nil
}

func UpdateLocation(id string, locationDataIn *dto.LocationInApp) (bool, error) {
	locationEntity := entity.NewLocationInApp(locationDataIn)
	locationOutDB := locationEntity.NewLocationOutDB()

	res, err := repository.UpdateOne(id, locationOutDB)
	if err != nil {
		return false, err
	}

	return res, nil
}

func DeleteLocation(id string) (bool, error) {
	res, err := repository.DeleteOne(id)
	if err != nil {
		return false, err
	}
	return res, nil
}
