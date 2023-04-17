package gokafka



func (s *Subscriber) Close() error {

	if err := s.consumer.Close(); err != nil {
		return err
	}
	
	return nil
}