package repository

import (
	"context"
	"database/sql"

	"kata-peya/internal/pet"
	"kata-peya/internal/tracer"

	sq "github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/codes"
)

type PetMysqlRepository struct {
	db *sql.DB
}

func NewPetMysqlRepository(db *sql.DB) *PetMysqlRepository {
	return &PetMysqlRepository{db: db}
}

func qB() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Question)
}

func (pr *PetMysqlRepository) QueryAll(ctx context.Context) ([]pet.Pet, error) {
	ctx, span := tracer.Start(ctx, "pets.repository.query_all")
	defer span.End()

	rows, err := qB().
		RunWith(pr.db).
		Select("*").
		From("pets").
		QueryContext(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	pets := make([]pet.Pet, 0)
	for rows.Next() {
		var p = pet.Pet{}
		if err := rows.Scan(&p.Id, &p.Name, &p.Vaccines, &p.AgeMonths); err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			return nil, err
		}

		pets = append(pets, p)
	}

	return pets, nil
}
