package api

import "context"

type TokenStore interface {
	Get(ctx context.Context) (string, error)
	Refresh(ctx context.Context) (string, error)
}
