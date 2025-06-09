package entity

import (
	"github.com/allansbo/goapi/internal/app/server/dto"
	"time"
)

type Coordinates struct {
	Latitude  string `bson:"latitude" json:"latitude"`
	Longitude string `bson:"longitude" json:"longitude"`
}

type Location struct {
	ID        string       `bson:"_id,omitempty" json:"id"`
	VehicleId string       `bson:"vehicle_id" json:"vehicle_id"`
	Timestamp time.Time    `bson:"timestamp" json:"timestamp"`
	Location  *Coordinates `bson:"location" json:"location"`
	Speed     int          `bson:"speed" json:"speed"`
	Status    string       `bson:"status" json:"status"`
}

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
