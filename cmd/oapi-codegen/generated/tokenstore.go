package api

import "context"

type TokenStore interface {
	Get(ctx context.Context) (string, error)
	Set(ctx context.Context, token string) error
	Refresh(ctx context.Context) (string, error)
	Revoke(ctx context.Context) error
	Valid(ctx context.Context) (bool, error)
}
