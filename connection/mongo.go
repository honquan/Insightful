package connection

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"insightful/common/config"
)

func InitMongo(conf *config.Config) *mongo.Database {

	ops := options.Client().ApplyURI(conf.Mongo.MongoURI).SetMaxPoolSize(500000)
	if conf.Mongo.MongoDebug {
		ops = ops.SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				zap.S().Debug(startedEvent.Command)
			},
		})
	}
	client, err := mongo.NewClient(ops)
	if err != nil {
		zap.S().Errorf("Error when create client mongo")
		return nil
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		zap.S().Errorf("Error when connect client mongo")
		return nil
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		zap.S().Errorf("Error when ping mongo")
		return nil
	}
	return client.Database(conf.Mongo.MongoDatabaseName)
}
