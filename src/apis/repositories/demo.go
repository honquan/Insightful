package repository

import (
	"gorm.io/gorm"
	"insightful/model"
)

type DemoRepository interface {
	Repository[*model.Demo]
}

type demoRepository struct {
	repository[*model.Demo]
}

func NewDemoRepository(db *gorm.DB) (DemoRepository, error) {
	return &demoRepository{
		repository: repository[*model.Demo]{
			db: db,
		},
	}, nil
}
