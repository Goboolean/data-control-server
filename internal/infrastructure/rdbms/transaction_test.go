package rdbms_test

import (
	"context"
	"testing"
	"database/sql"

)




func Test_Commit(t *testing.T) {

	count, err := queries.CountTestTableEntity(context.Background())
	if err != nil {
		t.Errorf("CountTestTableEntity() failed: %v", err)
		return
	}

	tx, err := db.NewTx(context.Background())
	if err != nil {
		t.Errorf("NewTx() failed: %v", err)
		return
	}

	q := queries.WithTx(tx.Transaction().(*sql.Tx))

	if err := q.InsertTestTableEntity(context.Background()); err != nil {
		t.Errorf("InsertTestTableEntity() failed: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		t.Errorf("Commit() failed: %v", err)
		return
	}

	updatedCount, err := queries.CountTestTableEntity(context.Background())
	if err != nil {
		t.Errorf("CountTestTableEntity() failed: %v", err)
		return
	}

	if updatedCount != count + 1 {
		t.Errorf("count = %d, updatedCount = %d, Commit() does not works", count, updatedCount)
		return
	}
}



func Test_Rollback(t *testing.T) {

	count, err := queries.CountTestTableEntity(context.Background())
	if err != nil {
		t.Errorf("CountTestTableEntity() failed: %v", err)
		return
	}

	tx, err := db.NewTx(context.Background())
	if err != nil {
		t.Errorf("NewTx() failed: %v", err)
		return
	}

	q := queries.WithTx(tx.Transaction().(*sql.Tx))

	if err := q.InsertTestTableEntity(context.Background()); err != nil {
		t.Errorf("InsertTestTableEntity() failed: %v", err)
		return
	}

	if err := tx.Rollback(); err != nil {
		t.Errorf("Rollback() failed: %v", err)
		return
	}

	updatedCount, err := queries.CountTestTableEntity(context.Background())
	if err != nil {
		t.Errorf("CountTestTableEntity() failed: %v", err)
		return
	}

	if updatedCount != count {
		t.Errorf("count = %d, updatedCount = %d, Rollback() does not works", count, updatedCount)
		return
	}
}