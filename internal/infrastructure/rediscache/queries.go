package rediscache



type Queries struct {
	rds *Redis
}

func New() *Queries {
	return &Queries{}
}


