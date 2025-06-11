package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// CoordinatesOutDB is the output data for saving a location in the database.
type CoordinatesOutDB struct {
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

// LocationOutDB is the output data for saving a location in the database.
type LocationOutDB struct {
	ID        string            `bson:"_id,omitempty"`
	VehicleId string            `bson:"vehicle_id"`
	Timestamp time.Time         `bson:"timestamp"`
	Location  *CoordinatesOutDB `bson:"location"`
	Speed     int               `bson:"speed"`
	Status    string            `bson:"status"`
}

// CoordinatesInDB is the input data for retrieving a location from the database.
type CoordinatesInDB struct {
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

// LocationInDB is the input data for retrieving a location from the database.
type LocationInDB struct {
	ID        bson.ObjectID    `bson:"_id"`
	VehicleId string           `bson:"vehicle_id"`
	Timestamp time.Time        `bson:"timestamp"`
	Location  *CoordinatesInDB `bson:"location"`
	Speed     int              `bson:"speed"`
	Status    string           `bson:"status"`
}

// QueryLocationOutDB is the input data for querying locations from the database.
type QueryLocationOutDB struct {
	Limit     int    `bson:"limit"`
	Page      int    `bson:"page"`
	VehicleId string `bson:"vehicle_id"`
	Status    string `bson:"status"`
}

// QueryLocationInDB is the input data for retrieving locations from the database.
type QueryLocationInDB struct {
	Limit int             `bson:"limit"`
	Page  int             `bson:"page"`
	Data  []*LocationInDB `bson:"data"`
}
