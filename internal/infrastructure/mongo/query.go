package mongo

import (
	"github.com/Goboolean/shared/pkg/resolver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Queries struct {
	db *DB
}

func New(db *DB) *Queries {
	return &Queries{db: db}
}



func (q *Queries) InsertStockBatch(tx resolver.Transactioner, stock string, batch []*StockAggregate) error {

	coll := q.db.client.Database(q.db.DefaultDatabase).Collection(stock)
	session := tx.Transaction().(mongo.Session)

	docs := make([]interface{}, len(batch))

	for idx := range batch {
		docs[idx] = &batch[idx]
	}

	return mongo.WithSession(tx.Context(), session, func(ctx mongo.SessionContext) error {
		_, err := coll.InsertMany(ctx, docs)
		return err
	})
}



func (q *Queries) FetchAllStockBatch(tx resolver.Transactioner, stock string) ([]*StockAggregate, error) {
	results := make([]*StockAggregate, 0)

	coll := q.db.client.Database(q.db.DefaultDatabase).Collection(stock)
	session := tx.Transaction().(mongo.Session)

	return results, mongo.WithSession(tx.Context(), session, func(ctx mongo.SessionContext) error {
		cursor, err := coll.Find(tx.Context(), bson.M{})
		if err != nil {
			return err
		}

		for cursor.Next(ctx) {
			var data *StockAggregate = &StockAggregate{}
			if err := cursor.Decode(data); err != nil {
				return err
			}

			results = append(results, data)
		}

		return cursor.Close(tx.Context())
	})
}



func (q *Queries) FetchAllStockBatchMassive(tx resolver.Transactioner, stock string, stockChan chan<- *StockAggregate) error {

	coll := q.db.client.Database(q.db.DefaultDatabase).Collection(stock)
	session := tx.Transaction().(mongo.Session)

	return mongo.WithSession(tx.Context(), session, func(ctx mongo.SessionContext) error {
		cursor, err := coll.Find(tx.Context(), bson.M{})
		if err != nil {
			return err
		}

		for cursor.Next(ctx) {
			var data *StockAggregate
			if err := cursor.Decode(&data); err != nil {
				return err
			}

			stockChan <- data
		}

		return cursor.Close(tx.Context())
	})
}



func (q *Queries) ClearAllStockData(tx resolver.Transactioner, stock string) error {

	coll := q.db.client.Database(q.db.DefaultDatabase).Collection(stock)
	session := tx.Transaction().(mongo.Session)

	return mongo.WithSession(tx.Context(), session, func(ctx mongo.SessionContext) error {
		_, err := coll.DeleteMany(ctx, bson.D{})
		return err
	})
}



func (q *Queries) GetStockDataLength(tx resolver.Transactioner, stock string) (length int, err error) {

	coll := q.db.client.Database(q.db.DefaultDatabase).Collection(stock)
	session := tx.Transaction().(mongo.Session)
	
	return length, mongo.WithSession(tx.Context(), session, func(ctx mongo.SessionContext) error {
		count, err := coll.CountDocuments(ctx, bson.D{})
		length = int(count)
		return err
	})
}