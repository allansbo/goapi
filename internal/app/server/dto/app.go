package dto

import "time"

// LocationInApp is the input data for the location endpoints
// that will be used to create or update a new location.
type LocationInApp struct {
	VehicleId string `validate:"required,alphanum,len=7" json:"vehicle_id"`
	Latitude  string `validate:"required,latitude" json:"latitude"`
	Longitude string `validate:"required,longitude" json:"longitude"`
	Status    string `validate:"required,oneof=moving stopped offline" json:"status"`
	Speed     int    `validate:"gte=0" json:"speed"`
}

// CoordinatesOutApp is the output data for the location endpoints
// that will be used to return a location.
type CoordinatesOutApp struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// LocationOutApp is the output data for the location endpoints
// that will be used to return a location.
type LocationOutApp struct {
	ID        string             `json:"id"`
	VehicleId string             `json:"vehicle_id"`
	Timestamp time.Time          `json:"timestamp"`
	Location  *CoordinatesOutApp `json:"location"`
	Speed     int                `json:"speed"`
	Status    string             `json:"status"`
}
