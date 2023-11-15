package repository

import (
	"context"
	"gorm.io/gorm"
	"insightful/model"
)

const txKey = "tx"

type Repository[T model.SQLModel] interface {
	DB(ctx context.Context) *gorm.DB
	BeginTx(ctx context.Context, fn func(tx context.Context) error) error
	Create(ctx context.Context, data T) error
	Save(ctx context.Context, data T) error
	Update(ctx context.Context, data T) error
	Delete(ctx context.Context, data T) error
	FindByID(ctx context.Context, id uint, preloads ...string) (T, error)
}

type repository[T model.SQLModel] struct {
	db *gorm.DB
}

func (r *repository[T]) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return r.db
}

func (r *repository[T]) SetDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, db)
}

func (r *repository[T]) BeginTx(ctx context.Context, fn func(tx context.Context) error) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = r.SetDB(ctx, tx)
		return fn(ctx)
	})
}

func (r *repository[T]) Create(ctx context.Context, data T) error {
	return r.DB(ctx).Model(data).Create(data).Error
}

func (r *repository[T]) Save(ctx context.Context, data T) error {
	return r.DB(ctx).Model(data).Save(data).Error
}

func (r *repository[T]) Update(ctx context.Context, data T) error {
	return r.DB(ctx).Model(data).Updates(data).Error
}

func (r *repository[T]) Delete(ctx context.Context, data T) error {
	return r.DB(ctx).Model(data).Delete(data).Error
}

func (r *repository[T]) FindByID(ctx context.Context, id uint, preloads ...string) (T, error) {
	var data T
	db := r.DB(ctx).Model(data)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if err := db.Where("id = ?", id).First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
