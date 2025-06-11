package db

import (
	"context"
	"fmt"
	"github.com/allansbo/goapi/internal/app/server/dto"
	"github.com/allansbo/goapi/internal/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoDBRepository implements the Repository interface for MongoDB operations.
type MongoDBRepository struct {
	client       *mongo.Client
	ctx          context.Context
	cancel       context.CancelFunc
	uri          string
	dbName       string
	dbCollection string
}

// NewMongoDBRepository creates a new instance of MongoDBRepository with the provided configuration.
func NewMongoDBRepository(cfg *config.EnvConfig) *MongoDBRepository {
	ctx, cancel := context.WithCancel(context.Background())

	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?retryWrites=true&w=majority&authSource=admin&ssl=false",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	clientOpts := options.Client().ApplyURI(uri)
	client, _ := mongo.Connect(clientOpts)

	return &MongoDBRepository{
		client:       client,
		dbName:       cfg.DBName,
		dbCollection: cfg.DBCollection,
		uri:          uri,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (m *MongoDBRepository) Stop() {
	m.cancel()
}

func (m *MongoDBRepository) Ping() error {
	if err := m.client.Ping(m.ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("mongodb ping failed: %w", err)
	}
	return nil
}

func (m *MongoDBRepository) collection() *mongo.Collection {
	return m.client.Database(m.dbName).Collection(m.dbCollection)
}

// InsertOne inserts a document into the collection.
func (m *MongoDBRepository) InsertOne(location *dto.LocationOutDB) (string, error) {
	res, err := m.collection().InsertOne(m.ctx, location)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(bson.ObjectID).Hex()
	return id, nil
}

// GetOne retrieves a single document by its ID from the collection.
func (m *MongoDBRepository) GetOne(id string) (*dto.LocationInDB, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var res bson.M
	err = m.collection().FindOne(m.ctx, bson.M{"_id": objectID}).Decode(&res)
	if err != nil {
		return nil, err
	}

	locationInDb := &dto.LocationInDB{}
	resBytes, err := bson.Marshal(res)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(resBytes, locationInDb); err != nil {
		return nil, err
	}

	return locationInDb, nil
}

// GetAll retrieves all documents from the collection, limited
// by the specified count  and filtered by the provided filter.
func (m *MongoDBRepository) GetAll(query *dto.QueryLocationOutDB) (*dto.QueryLocationInDB, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	filter := bson.M{}
	if query.VehicleId != "" {
		filter["vehicle_id"] = query.VehicleId
	}
	if query.Status != "" {
		filter["status"] = query.Status
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((query.Page - 1) * query.Limit)).SetLimit(int64(query.Limit))

	cursor, err := m.collection().Find(m.ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var locations []*dto.LocationInDB
	if err := cursor.All(m.ctx, &locations); err != nil {
		return nil, err
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	qLocationsInDB := new(dto.QueryLocationInDB)
	qLocationsInDB.Limit = query.Limit
	qLocationsInDB.Page = query.Page
	qLocationsInDB.Data = locations

	return qLocationsInDB, nil
}

// UpdateOne updates a single document by its ID in the collection.
func (m *MongoDBRepository) UpdateOne(id string, location *dto.LocationOutDB) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{"$set": location}

	res, err := m.collection().UpdateOne(m.ctx, bson.M{"_id": objectID}, data)
	if err != nil {
		return false, err
	}

	return res.ModifiedCount == 1, nil
}

// DeleteOne deletes a single document by its ID from the collection.
func (m *MongoDBRepository) DeleteOne(id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	res, err := m.collection().DeleteOne(m.ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	return res.DeletedCount == 1, nil
}
