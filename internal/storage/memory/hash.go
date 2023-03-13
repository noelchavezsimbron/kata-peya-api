package memory

import (
	"context"

	"kata-peya/internal/tracer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	key interface {
		string | int | int8 | int16 | int32 | int64
	}

	HashDb[K key, T any] struct {
		m map[K]T
	}
)

func NewHashDB[K key, T any](data map[K]T) *HashDb[K, T] {
	db := &HashDb[K, T]{m: map[K]T{}}
	if len(data) > 0 {
		db.m = data
	}
	return db
}

func (hdb *HashDb[K, T]) Get(ctx context.Context, key K) T {
	ctx, span := tracer.Start(ctx, "HashDb.Get")
	setSpanDBTags(span)
	defer span.End()

	return hdb.m[key]
}

func (hdb *HashDb[K, T]) Set(ctx context.Context, key K, e T) {
	ctx, span := tracer.Start(ctx, "HashDb.Set")
	setSpanDBTags(span)
	defer span.End()

	hdb.m[key] = e
}

func (hdb *HashDb[K, T]) List(ctx context.Context) []T {
	ctx, span := tracer.Start(ctx, "HashDb.List")
	setSpanDBTags(span)
	defer span.End()

	r := make([]T, 0, len(hdb.m))
	for _, v := range hdb.m {
		r = append(r, v)
	}
	return r
}

func setSpanDBTags(span trace.Span) {
	span.SetAttributes(attribute.String("db.system", "hash"))
	span.SetAttributes(attribute.String("db.connection_string", "local"))
}
