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
	"log/slog"
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

	slog.Info("mongodb connection established")
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
	slog.Info("document inserted", "id", id)
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
