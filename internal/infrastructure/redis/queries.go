package redis

type Queries struct {
	rds *Redis
}

func New() *Queries {
	return &Queries{}
}
