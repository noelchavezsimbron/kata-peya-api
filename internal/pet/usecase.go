package pet

import "context"

type UseCase interface {
	FindAll(ctx context.Context) ([]Pet, error)
}
