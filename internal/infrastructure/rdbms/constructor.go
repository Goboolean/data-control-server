package rdbms

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Goboolean/shared/pkg/resolver"
	_ "github.com/lib/pq"
)

type PSQL struct {
	db *sql.DB
}

func NewDB(c *resolver.ConfigMap) (*PSQL, error) {
	user, err := c.GetStringKey("USER")
	if err != nil {
		return nil, err
	}

	password, err := c.GetStringKey("PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	database, err := c.GetStringKey("DATABASE")
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	return &PSQL{
		db: db,
	}, nil
}

func (p *PSQL) Close() error {
	return p.db.Close()
}

func (p *PSQL) Ping() error {
	return p.db.Ping()
}

func (p *PSQL) NewTx(ctx context.Context) (resolver.Transactioner, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	return NewTransaction(tx, ctx), err
}

func NewQueries(db *PSQL) *Queries {
	return &Queries{db: db.db}
}
