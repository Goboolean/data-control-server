package mongodb

import (
	"context"
	"fmt"
	"os"

	adapter "github.com/Goboolean/data-control-server/internal/adaptor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client          *mongo.Client
	defaultDatabase string
}

var (
	MONGO_HOST     = os.Getenv("MONGO_HOST")
	MONGO_PORT     = os.Getenv("MONGO_PORT")
	MONGO_USER     = os.Getenv("MONGO_USER")
	MONGO_PASS     = os.Getenv("MONGO_PASS")
	MONGO_DATABASE = os.Getenv("MONGO_DATABASE")
)

var mongoInstance *Mongo

func NewMongo() *Mongo {
	if mongoInstance == nil {

		mongoInstance = &Mongo{defaultDatabase: os.Getenv("MONGO_DATABASE")}

		var mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority",
			MONGO_USER, MONGO_PASS, MONGO_HOST, MONGO_PORT)
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

		client, err := mongo.Connect(context.TODO(), opts)

		if err != nil {
			panic(err)
		}

		mongoInstance.client = client
	}

	return mongoInstance
}

func (m *Mongo) InsertOne(transaction adapter.Transaction, data *interface{}, collName string) error {
	coll := m.client.Database(m.defaultDatabase).Collection(collName)
	_, err := coll.InsertOne(transaction.Context(), data)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mongo) InsertMany(transaction adapter.Transaction, data []interface{}, collName string) error {
	coll := m.client.Database(m.defaultDatabase).Collection(collName)

	_, err := coll.InsertMany(transaction.Context(), data)

	if err != nil {
		return err
	}

	return nil
}
