package model

import (
	"gorm.io/gorm"
)

type Demo struct {
	gorm.Model
}

func (Demo) TableName() string {
	return "demo"
}
