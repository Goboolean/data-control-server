package mongodb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var (
	MONGO_HOST     = os.Getenv("MONGO_HOST")
	MONGO_PORT     = os.Getenv("MONGO_PORT")
	MONGO_USER     = os.Getenv("MONGO_USER")
	MONGO_PASS     = os.Getenv("MONGO_PASS")
	MONGO_DATABASE = os.Getenv("MONGO_DATABASE")
)

var instance *mongo.Client

func NewInstance() *mongo.Client {
	if instance == nil {

		var mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority",
			MONGO_USER, MONGO_PASS, MONGO_HOST, MONGO_PORT)
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

		client, err := mongo.Connect(context.TODO(), opts)

		if err != nil {
			panic(err)
		}

		instance = client
	}

	return instance
}

