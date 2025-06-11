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
// The data is coming from the database.
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

// QueryLocationRequest is the entity that represents a request to query locations.
type QueryLocationRequest struct {
	Limit     int    `bson:"limit" json:"limit"`
	Page      int    `bson:"page" json:"page"`
	VehicleId string `bson:"vehicle_id" json:"vehicle_id"`
	Status    string `bson:"status" json:"status"`
}

// NewQueryLocationRequest is a function that creates a new query location request
func NewQueryLocationRequest(query *dto.QueryLocationRequest) *QueryLocationRequest {
	return &QueryLocationRequest{
		Limit:     query.Limit,
		Page:      query.Page,
		VehicleId: query.VehicleId,
		Status:    query.Status,
	}
}

// NewQueryLocationOutDB is a function that exports the query location request to the database format.
func (q *QueryLocationRequest) NewQueryLocationOutDB() *dto.QueryLocationOutDB {
	return &dto.QueryLocationOutDB{
		Limit:     q.Limit,
		Page:      q.Page,
		VehicleId: q.VehicleId,
		Status:    q.Status,
	}
}

// PaginationInfo is the entity that represents pagination information for a query response.
type PaginationInfo struct {
	Page  int
	Limit int
}

// QueryLocationResponse is the entity that represents a response to a query for locations.
type QueryLocationResponse struct {
	Data       []*Location
	Pagination *PaginationInfo
}

// NewQueryLocationResponse is a function that creates a new query location response from a database query result.
func NewQueryLocationResponse(q *dto.QueryLocationInDB) *QueryLocationResponse {
	dataLocations := make([]*Location, 0, len(q.Data))
	for _, loc := range q.Data {
		locEntity := NewLocationInDB(loc)
		dataLocations = append(dataLocations, locEntity)
	}

	pageInfo := &PaginationInfo{
		Limit: q.Limit,
		Page:  q.Page,
	}
	return &QueryLocationResponse{
		Pagination: pageInfo,
		Data:       dataLocations,
	}
}

// NewQueryLocationOutApp is a function that exports the query location response to the user
func (q *QueryLocationResponse) NewQueryLocationOutApp() *dto.QueryLocationResponse {
	dataLocations := make([]*dto.LocationOutApp, 0, len(q.Data))
	for _, loc := range q.Data {
		locOutApp := loc.NewLocationOutApp()
		dataLocations = append(dataLocations, locOutApp)
	}

	return &dto.QueryLocationResponse{
		Success: len(dataLocations) != 0,
		Data:    dataLocations,
		Pagination: &dto.PaginationInfoResponse{
			Page:  q.Pagination.Page,
			Limit: q.Pagination.Limit,
		},
	}
}
