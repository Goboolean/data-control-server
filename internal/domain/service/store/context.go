package store

import "context"



type contextController struct {
	ctx context.Context
	cancel context.CancelFunc
}


func (c *contextController) Context() context.Context {
	return c.ctx
}

func (c *contextController) Done() <- chan struct{} {
	return c.ctx.Done()
}

func new_ctx(ctx context.Context) *contextController {
	ctx, cancel := context.WithCancel(ctx)
	return &contextController{ctx: ctx, cancel: cancel}
}