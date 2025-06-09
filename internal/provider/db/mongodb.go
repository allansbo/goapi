package db

import (
	"context"
	"fmt"
	"github.com/allansbo/goapi/internal/app/server/dto"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log/slog"
	"time"

	"github.com/allansbo/goapi/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDBRepository struct {
	DBUser       string
	DBPass       string
	DBName       string
	DBCollection string
	DBHost       string
	DBPort       string
}

func NewMongoDBRepository(cfg *config.EnvConfig) *MongoDBRepository {
	return &MongoDBRepository{
		DBUser:       cfg.DBUser,
		DBPass:       cfg.DBPass,
		DBName:       cfg.DBName,
		DBCollection: cfg.DBCollection,
		DBHost:       cfg.DBHost,
		DBPort:       cfg.DBPort,
	}
}

func (m *MongoDBRepository) getMongoDBURI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?retryWrites=true&w=majority&authSource=admin&ssl=false",
		m.DBUser, m.DBPass, m.DBHost, m.DBPort, m.DBName,
	)
}

// Connect create a connection to the mongodb database and use Ping to check if the connection is established.
func (m *MongoDBRepository) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(m.getMongoDBURI())
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			slog.Error("error on disconnect from mongodb", "error", err.Error())
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	slog.Info("mongodb connection established")

	return nil
}

// InsertOne inserts a document into the collection
func (m *MongoDBRepository) InsertOne(location *dto.LocationOutDB) (string, error) {
	data, err := bson.Marshal(location)
	if err != nil {
		return "", err
	}

	clientOpts := options.Client().ApplyURI(m.getMongoDBURI())
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			slog.Error("error on disconnect from mongodb", "error", err.Error())
		}
	}()

	collection := client.Database(m.DBName).Collection(m.DBCollection)
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return "", err
	}

	slog.Info("document inserted successfully", "id", res.InsertedID, "document", location)

	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (m *MongoDBRepository) GetOne(id string) (*dto.LocationInDB, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	clientOpts := options.Client().ApplyURI(m.getMongoDBURI())
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			slog.Error("error on disconnect from mongodb", "error", err.Error())
		}
	}()

	collection := client.Database(m.DBName).Collection(m.DBCollection)
	var res bson.M
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&res)
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

func (m *MongoDBRepository) UpdateOne(id string, location *dto.LocationOutDB) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	updateData := map[string]interface{}{"$set": location}

	data, err := bson.Marshal(updateData)
	if err != nil {
		return false, err
	}

	clientOpts := options.Client().ApplyURI(m.getMongoDBURI())
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return false, fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			slog.Error("error on disconnect from mongodb", "error", err.Error())
		}
	}()

	collection := client.Database(m.DBName).Collection(m.DBCollection)
	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, data)
	if err != nil {
		return false, err
	}

	return res.ModifiedCount == 1, nil
}

func (m *MongoDBRepository) DeleteOne(id string) (bool, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	clientOpts := options.Client().ApplyURI(m.getMongoDBURI())
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return false, fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			slog.Error("error on disconnect from mongodb", "error", err.Error())
		}
	}()

	collection := client.Database(m.DBName).Collection(m.DBCollection)
	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}

	return res.DeletedCount == 1, nil
}
