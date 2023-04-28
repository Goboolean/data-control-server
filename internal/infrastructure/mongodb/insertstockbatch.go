package mongodb

import (
	infra "github.com/Goboolean/stock-fetch-server/internal/infrastructure/transaction"
	"go.mongodb.org/mongo-driver/mongo"
)

func (q *Queries) InsertStockBatch(tx infra.Transactioner, stock string, batch []StockAggregate) error {

	coll := q.client.Database(MONGO_DATABASE).Collection(stock)
	session := tx.Transaction().(mongo.Session)

	docs := make([]interface{}, len(batch))

	for idx := range batch {
		docs[idx] = &batch[idx]
	}

	_, err := session.WithTransaction(tx.Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		return coll.InsertMany(ctx, docs)
	})

	return err
}
