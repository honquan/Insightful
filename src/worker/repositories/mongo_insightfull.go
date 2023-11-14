package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"insightful/common/pagination"
	"insightful/model"
)

type MongoInsightfullRepository interface {
	Repository
	FindMany(ctx context.Context, p *pagination.Pagination, filter bson.M) ([]*model.Insightful, error)
}

type mongoInsightfullRepository struct {
	mongoDBRepo
}

func NewMongoInsightfullRepository(db *mongo.Database) MongoInsightfullRepository {
	return &mongoInsightfullRepository{
		mongoDBRepo{db},
	}
}

func (r *mongoInsightfullRepository) FindMany(ctx context.Context, p *pagination.Pagination, filter bson.M) ([]*model.Insightful, error) {
	var (
		err        error
		insightful []*model.Insightful
	)
	c := r.Collection(ctx, model.Insightful{}, -1)
	if !p.NoUse {
		p.Correct()
		total, err := c.CountDocuments(ctx, filter)
		if err != nil {
			return nil, err
		}
		p.SetTotal(total)
	}
	it, err := c.Find(ctx, filter, p.FindOptionWithoutSkip().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		return nil, err
	}
	return insightful, it.All(ctx, &insightful)
}
