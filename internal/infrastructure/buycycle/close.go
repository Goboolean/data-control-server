package buycycle



func (s *Subscriber) Close() error {

	if err := s.conn.Close(); err != nil {
		return err
	}
	
	return nil
}