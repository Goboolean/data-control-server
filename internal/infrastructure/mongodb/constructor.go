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
	MONGO_URI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority",
	MONGO_USER, MONGO_PASS, MONGO_HOST, MONGO_PORT)
)

var instance *mongo.Client

func init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	instance = client
}



func NewInstance() *mongo.Client {
	return instance
}

func Close() error {

	if err := instance.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}