package dto

import "time"

type LocationInApp struct {
	VehicleId string `validate:"required,alphanum,len=7" json:"vehicle_id"`
	Latitude  string `validate:"required,latitude" json:"latitude"`
	Longitude string `validate:"required,longitude" json:"longitude"`
	Status    string `validate:"required,oneof=moving stopped offline" json:"status"`
	Speed     int    `validate:"gte=0" json:"speed"`
}

type CoordinatesOutApp struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type LocationOutApp struct {
	ID        string             `json:"id"`
	VehicleId string             `json:"vehicle_id"`
	Timestamp time.Time          `json:"timestamp"`
	Location  *CoordinatesOutApp `json:"location"`
	Speed     int                `json:"speed"`
	Status    string             `json:"status"`
}
