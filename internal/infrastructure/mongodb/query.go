package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type Queries struct {
	client *mongo.Client
	tx *mongo.Session
}

func New() *Queries {
	return &Queries{client: NewInstance()}
}


