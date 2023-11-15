package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"insightful/model"
)

type InsightfullRepository interface {
	MongoRepository
	Create(ctx context.Context, data *model.Insightful) error
	CreateMany(ctx context.Context, data []interface{}) error
	BulkWrite(ctx context.Context, data []mongo.WriteModel) error
}

type insightfullRepository struct {
	mongoDBRepo
}

func NewInsightfullRepository(db *mongo.Database) InsightfullRepository {
	return &insightfullRepository{mongoDBRepo{db}}
}

func (r *insightfullRepository) Create(ctx context.Context, data *model.Insightful) error {
	_, err := r.Collection(ctx, model.Insightful{}).InsertOne(ctx, data)

	return err
}

func (r *insightfullRepository) CreateMany(ctx context.Context, data []interface{}) error {
	_, err := r.Collection(ctx, model.Insightful{}).InsertMany(ctx, data)

	return err
}

func (r *insightfullRepository) BulkWrite(ctx context.Context, data []mongo.WriteModel) error {
	// begin bulk
	coll := r.Collection(ctx, model.Insightful{})
	opts := options.BulkWrite().SetOrdered(false)

	_, err := coll.BulkWrite(context.TODO(), data, opts)
	if err != nil {
		return err
	}

	return nil
}
