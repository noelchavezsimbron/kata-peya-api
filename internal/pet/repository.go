package pet

import "context"

type Repository interface {
	QueryAll(ctx context.Context) ([]Pet, error)
}
