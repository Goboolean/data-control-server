package store

import (
	"context"
	"fmt"
)




type Store struct {
	c *contextController

	m map[string] * contextController
}


func New(ctx context.Context) *Store {
	return &Store{
		c: new_ctx(ctx),
		m: make(map[string] * contextController),
	}
}


func (s *Store) StockExists(stock string) bool {
	_, ok := s.m[stock]
	return ok
}


func (s *Store) StoreStock(stock string) error {
	if s.StockExists(stock) {
		return fmt.Errorf("stock %s already exists", stock)
	}

	s.m[stock] = new_ctx(s.c.ctx)
	return nil
}


func (s *Store) UnstoreStock(stock string) error {
	if !s.StockExists(stock) {
		return fmt.Errorf("stock %s not exists", stock)
	}
	s.m[stock].cancel()
	delete(s.m, stock)
	return nil
}


func (s *Store) Close() {
	s.m = nil
	s.c.cancel()
}