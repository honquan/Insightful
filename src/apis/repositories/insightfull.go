package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"insightful/model"
)

type InsightfullRepository interface {
	Repository
	Create(ctx context.Context, data *model.Insightful) error
	CreateMany(ctx context.Context, data []model.Insightful) error
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

func (r *insightfullRepository) CreateMany(ctx context.Context, data []model.Insightful) error {
	arrInsightfull := []interface{}{}
	for _, insightfull := range data {
		arrInsightfull = append(arrInsightfull, insightfull)
	}
	_, err := r.Collection(ctx, model.Insightful{}).InsertMany(ctx, arrInsightfull)

	return err
}
