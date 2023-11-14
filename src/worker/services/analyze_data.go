package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"insightful/common/config"
	"insightful/common/pagination"
	"insightful/src/worker/dtos"
	"insightful/src/worker/repositories"
)

type AnalyzeDataService interface {
	TriggerAnalyzeData(ctx context.Context, args []string) error
}

type analyzeDataService struct {
	config                  *config.Config
	postgresInsightfullRepo repositories.PostgresInsightfullRepository
	mongoInsightfullRepo    repositories.MongoInsightfullRepository
}

func NewAnalyzeDataService(
	cfg *config.Config,
	postgresInsightfullRepo repositories.PostgresInsightfullRepository,
	mongoInsightfullRepo repositories.MongoInsightfullRepository,
) AnalyzeDataService {
	return &analyzeDataService{
		config:                  cfg,
		postgresInsightfullRepo: postgresInsightfullRepo,
		mongoInsightfullRepo:    mongoInsightfullRepo,
	}
}

func (s *analyzeDataService) TriggerAnalyzeData(ctx context.Context, args []string) error {
	// get data by create at
	p := &pagination.Pagination{
		Limit: dtos.LimitEachQuery,
	}
	andFilter, err := dtos.QueryForDistributeInsightful(args)
	if err != nil {
		zap.S().Errorw("Error when QueryForDistributeInsightful", "error", err)
		return err
	}

	for {
		data, err := s.mongoInsightfullRepo.FindMany(ctx, p, bson.M{"$and": andFilter})
		if err != nil {
			zap.S().Errorw("error when call mongoInsightfullRepo.FindMany", "error", err, "pagination", p, "query", andFilter)
			return err
		}
		if data == nil || len(data) == 0 {
			break
		}

		// TODO process add to postgres

		// get last in data
		lastData := data[len(data)-1]
		if lastData != nil {
			andFilter = append(andFilter, bson.M{
				"created_at": bson.M{"$gt": lastData.CreatedAt},
			})
		}

		// check last
		if len(data) < dtos.LimitEachQuery {
			break
		}
	}

	return nil
}
