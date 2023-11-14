package dtos

import (
	"go.mongodb.org/mongo-driver/bson"
	"insightful/src/worker/utils"
)

const LimitEachQuery = 10

func QueryForDistributeInsightful(args []string) ([]bson.M, error) {
	andFilter := []bson.M{
		{"updated_at": bson.M{"$exists": false}},
	}

	// check create at for distribute worker
	startIn, endIn, err := utils.DistributeTimeByteArgument(args)
	if err != nil {
		return nil, err
	}
	if startIn > 0 {
		andFilter = append(andFilter, bson.M{
			"created_at": bson.M{"$gte": startIn},
		})
	}

	if endIn > 0 {
		andFilter = append(andFilter, bson.M{
			"created_at": bson.M{"$lt": endIn},
		})
	}

	return andFilter, nil
}
