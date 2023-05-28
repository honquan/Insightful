package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Insightful struct {
	Mongo     `bson:",inline"`
	Done      byte        `bson:"done" json:"done"`
	Coodiates interface{} `bson:"coodiates" json:"coodiates,omitempty"`
}

func (Insightful) CollectionName(ws string) string {
	return fmt.Sprintf("%v._insightfull", ws)
}

func (c *Insightful) UpdateModel(ctx context.Context) bson.M {
	return bson.M{
		"done":       1,
		"updated_at": time.Now(),
	}
}
