// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package rdbms

import (
	"database/sql"
	"time"
)

type Platform struct {
	Name        string
	Description sql.NullString
}

type ProductMetum struct {
	ID          string
	Name        string
	Symbol      string
	Description sql.NullString
	Type        string
	Exchange    string
	Location    sql.NullString
}

type ProductPlatform struct {
	ProductID    string
	PlatformName string
	Identifier   string
}

type StoreLog struct {
	ProductID string
	StoredAt  time.Time
	Status    string
}

type TestTable struct {
	ID int32
}
