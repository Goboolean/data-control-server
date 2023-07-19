package store

import "context"



type contextController struct {
	ctx context.Context
	cancel context.CancelFunc
}

func new_ctx(ctx context.Context) *contextController {
	ctx, cancel := context.WithCancel(ctx)
	return &contextController{ctx: ctx, cancel: cancel}
}