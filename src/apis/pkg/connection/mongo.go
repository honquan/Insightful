package connection

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"insightful/src/apis/conf"
)

func NewMongoConnection() (*mongo.Database, func() error, error) {

	ops := options.Client().ApplyURI(conf.EnvConfig.MongoURI)
	if conf.EnvConfig.MongoDebug {
		ops = ops.SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				zap.S().Debug(startedEvent.Command)
			},
		})
	}
	client, err := mongo.NewClient(ops)
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, nil, err
	}
	return client.Database(conf.EnvConfig.MongoDatabaseName), func() error {
		return client.Disconnect(ctx)
	}, nil
}
