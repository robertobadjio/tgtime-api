package router

import (
	"context"
)

type Repository interface {
	CreateRouter(ctx context.Context, router *Router) (*Router, error)
	UpdateRouter(ctx context.Context, router *Router) (*Router, error)
	GetRouter(ctx context.Context, routerId int) (*Router, error)
	GetRouters(ctx context.Context) ([]*Router, error)
	DeleteRouter(ctx context.Context, routerId int) error
}
