package usecase

import (
	"context"

	pet2 "kata-peya/internal/pet"
	"kata-peya/internal/tracer"

	"go.opentelemetry.io/otel/codes"
)

type UseCase struct {
	repo pet2.Repository
}

func NewUseCase(repo pet2.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) FindAll(ctx context.Context) ([]pet2.Pet, error) {
	ctx, span := tracer.Start(ctx, "pets.use_case.find_all")
	defer span.End()

	pets, err := uc.repo.QueryAll(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	return pets, nil
}
