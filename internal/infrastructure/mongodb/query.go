package mongodb

import (
	"github.com/Goboolean/data-control-server/internal/infrastructure/transaction"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)



type Queries struct {
	client *mongo.Client
}

func NewQueries(client *mongo.Client) *Queries {
	return &Queries{client: client}
}

func (q *Queries) InsertStockList(tx infratx.TransactionHandler, collName string, arg []StockAggregate) error {

	coll := q.client.Database(MONGO_DATABASE).Collection(collName)

	var data []interface{}
	for idx := range arg {

		segment, err := bson.Marshal(arg[idx])

		if err != nil {
			return err
		}

		data = append(data, segment)
	}

	session := tx.Transaction().(mongo.Session)

	_, err := session.WithTransaction(tx.Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		return coll.InsertMany(ctx, data)
	})

	return err
}