package mongo

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type (
	Collection interface {
		Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)
		FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
		FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
		InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
		DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	}

	Database interface {
		Collection(name string) Collection
		Close()
	}

	Config struct {
		ConnectionURI  string
		ConnectTimeout int
		Database       string
	}

	documentDb struct {
		db *mongo.Database
	}
)

func New(cfg Config) (Database, error) {
	log.Printf("ConnectionURI: %s\n", cfg.ConnectionURI)

	opts := options.Client().
		ApplyURI(cfg.ConnectionURI).
		SetReadPreference(readpref.SecondaryPreferred()).
		SetTimeout(time.Duration(30) * time.Second).
		SetMonitor(otelmongo.NewMonitor())

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create documentDB client: %v", err)
	}

	timeout := time.Duration(cfg.ConnectTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to documentDB cluster: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping documentDB cluster: %v", err)
	}

	return &documentDb{db: client.Database(cfg.Database)}, nil
}

func (d *documentDb) Collection(name string) Collection {
	return d.db.Collection(name)
}

func (d *documentDb) Close() {
	if err := d.db.Client().Disconnect(context.TODO()); err != nil {
		log.Error(err)
	}
}
