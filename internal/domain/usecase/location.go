package usecase

import (
	"github.com/allansbo/goapi/internal/app/server/dto"
	"github.com/allansbo/goapi/internal/domain/entity"
	"github.com/allansbo/goapi/internal/provider/db"
)

type locationUseCase struct {
	repository db.Repository
}

var l locationUseCase

func LoadLocationUseCase(repository db.Repository) {
	l.repository = repository
}

// SaveLocation saves a new location in the database and returns the saved location.
// It takes a pointer to dto.LocationInApp as input, which contains the validated location data.
// It returns a pointer to dto.LocationOutApp and an error if any occurs.
func SaveLocation(locationDataIn *dto.LocationInApp) (*dto.LocationOutApp, error) {
	locationEntity := entity.NewLocationInApp(locationDataIn)
	locationOutDB := locationEntity.NewLocationOutDB()

	var err error
	locationEntity.ID, err = l.repository.InsertOne(locationOutDB)
	if err != nil {
		return nil, err
	}

	return locationEntity.NewLocationOutApp(), nil
}

// GetLocationById retrieves a location by its ID from the database.
// It takes a string ID as input and returns a pointer to dto.LocationOutApp and an error if any occurs.
func GetLocationById(id string) (*dto.LocationOutApp, error) {
	locationInDB, err := l.repository.GetOne(id)
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

	res, err := l.repository.UpdateOne(id, locationOutDB)
	if err != nil {
		return false, err
	}

	return res, nil
}

// DeleteLocation deletes a location by its ID from the database.
// It takes a string ID as input and returns a boolean indicating success and an error if any occurs.
func DeleteLocation(id string) (bool, error) {
	res, err := l.repository.DeleteOne(id)
	if err != nil {
		return false, err
	}
	return res, nil
}
