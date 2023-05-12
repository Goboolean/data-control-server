package relayer

import "fmt"



type store struct {
	m map[string]struct{}
}



func (s *store) stockExists(stock string) bool {
	_, ok := s.m[stock]
	return ok
}

func (s *store) storeStock(stock string) error {
	if s.stockExists(stock) {
		return fmt.Errorf("stock %s already exists", stock)
	}
	return nil
}

func (s *store) unstoreStock(stock string) error {
	if !s.stockExists(stock) {
		return fmt.Errorf("stock %s not exists", stock)
	}
	return nil
}