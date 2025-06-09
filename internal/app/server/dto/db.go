package dto

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type CoordinatesOutDB struct {
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type LocationOutDB struct {
	ID        string            `bson:"_id,omitempty"`
	VehicleId string            `bson:"vehicle_id"`
	Timestamp time.Time         `bson:"timestamp"`
	Location  *CoordinatesOutDB `bson:"location"`
	Speed     int               `bson:"speed"`
	Status    string            `bson:"status"`
}

type CoordinatesInDB struct {
	Latitude  string `bson:"latitude"`
	Longitude string `bson:"longitude"`
}

type LocationInDB struct {
	ID        bson.ObjectID    `bson:"_id"`
	VehicleId string           `bson:"vehicle_id"`
	Timestamp time.Time        `bson:"timestamp"`
	Location  *CoordinatesInDB `bson:"location"`
	Speed     int              `bson:"speed"`
	Status    string           `bson:"status"`
}
