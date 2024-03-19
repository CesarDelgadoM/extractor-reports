package database

import (
	"context"
	"fmt"

	"github.com/CesarDelgadoM/extractor-reports/config"
	"github.com/CesarDelgadoM/extractor-reports/pkg/logger/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	config *config.MongoConfig
}

func ConnectMongoDB(config *config.MongoConfig) *MongoDB {
	uri := fmt.Sprintf(config.URI, config.User, config.Password)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		zap.Log.Fatal(err)
	}

	zap.Log.Info("Connection to mongodb success")
	return &MongoDB{
		Client: client,
		config: config,
	}
}

func (m *MongoDB) Disconnect() error {
	return m.Client.Disconnect(context.TODO())
}

func (m *MongoDB) CollectionRestaurant() *mongo.Collection {
	return m.Client.Database(m.config.DBName).Collection("restaurant")
}
