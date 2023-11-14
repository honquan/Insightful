package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Insightful struct {
	Mongo       `bson:",inline"`
	Coordinates interface{} `bson:"c" json:"c,omitempty"`
	//Done      byte        `bson:"done" json:"done"`
}

func (Insightful) CollectionName(ws string) string {
	return fmt.Sprintf("%v._insightfull", ws)
}

func (c *Insightful) UpdateModel(ctx context.Context) bson.M {
	return bson.M{
		"updated_at": time.Now().Unix(),
	}
}
