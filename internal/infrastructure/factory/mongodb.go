package factory

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type Mongo struct {
	client *mongo.Client
	defaultDatabase string
}

var mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/?maxPoolSize=20&w=majority",
  os.Getenv("MONGO_USER"),
  os.Getenv("MONGO_PASS"),
  os.Getenv("MONGO_HOST"),
  os.Getenv("MONGO_PORT"),
)


var instance *Mongo

func NewMongo() *Mongo {
	if instance == nil {
		instance = &Mongo{defaultDatabase: os.Getenv("MONGO_DATABASE")}
		instance.Connect()
	}
	return instance
}

func (m *Mongo) Connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	m.client = client
}



func (m *Mongo) InsertOne(ctx context.Context, data interface{}, collName string) error {
	coll := m.client.Database(m.defaultDatabase).Collection(collName)

	_, err := coll.InsertOne(ctx, data)

	if err != nil {
		return err
	}

	return nil
}



func (m *Mongo) InsertMany(ctx context.Context, data []interface{}, collName string) error {
	coll := m.client.Database(m.defaultDatabase).Collection(collName)

	_, err := coll.InsertMany(ctx, data)

	if err != nil {
		return err
	}

	return nil
}