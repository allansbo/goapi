package entity

import (
	"time"

	"github.com/allansbo/goapi/internal/app/server/dto"
)

// Coordinates is the entity that represents the coordinates of a location.
type Coordinates struct {
	Latitude  string `bson:"latitude" json:"latitude"`
	Longitude string `bson:"longitude" json:"longitude"`
}

// Location is the entity that represents the location of a vehicle.
type Location struct {
	ID        string       `bson:"_id,omitempty" json:"id"`
	VehicleId string       `bson:"vehicle_id" json:"vehicle_id"`
	Timestamp time.Time    `bson:"timestamp" json:"timestamp"`
	Location  *Coordinates `bson:"location" json:"location"`
	Speed     int          `bson:"speed" json:"speed"`
	Status    string       `bson:"status" json:"status"`
}

// NewLocationInApp is a function that creates a new location in the application.
// The user input was validated by the *dto.LocationInApp struct.
func NewLocationInApp(location *dto.LocationInApp) *Location {
	return &Location{
		VehicleId: location.VehicleId,
		Timestamp: time.Now(),
		Speed:     location.Speed,
		Status:    location.Status,
		Location: &Coordinates{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
	}
}

// NewLocationInDB is a function that creates a new location in the application.
// The data is comming from the database.
func NewLocationInDB(location *dto.LocationInDB) *Location {
	return &Location{
		ID:        location.ID.Hex(),
		VehicleId: location.VehicleId,
		Timestamp: location.Timestamp,
		Speed:     location.Speed,
		Status:    location.Status,
		Location: &Coordinates{
			Latitude:  location.Location.Latitude,
			Longitude: location.Location.Longitude,
		},
	}
}

// NewLocationOutDB is a function that exports the location to the database format.
func (l *Location) NewLocationOutDB() *dto.LocationOutDB {
	return &dto.LocationOutDB{
		VehicleId: l.VehicleId,
		Timestamp: l.Timestamp,
		Speed:     l.Speed,
		Status:    l.Status,
		Location: &dto.CoordinatesOutDB{
			Latitude:  l.Location.Latitude,
			Longitude: l.Location.Longitude,
		},
	}
}

// NewLocationOutApp is a function that exports the location
// to the format that will response a request user.
func (l *Location) NewLocationOutApp() *dto.LocationOutApp {
	return &dto.LocationOutApp{
		ID:        l.ID,
		VehicleId: l.VehicleId,
		Timestamp: l.Timestamp,
		Speed:     l.Speed,
		Status:    l.Status,
		Location: &dto.CoordinatesOutApp{
			Latitude:  l.Location.Latitude,
			Longitude: l.Location.Longitude,
		},
	}
}
