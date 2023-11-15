package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	model "insightful/model"
	"time"
)

type MongoRepository interface {
	BeginTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type mongoColRepo struct {
	collection *mongo.Collection
}

type mongoDBRepo struct {
	db *mongo.Database
}

func (d mongoDBRepo) Collection(ctx context.Context, i model.MongoModel) *mongo.Collection {
	date := time.Now().Format("2006-01-02")

	colName := i.CollectionName(date)
	return d.db.Collection(colName)
}

func (d mongoDBRepo) NewSession() (mongo.Session, error) {
	session, err := d.db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	if err := session.StartTransaction(); err != nil {
		return nil, err
	}
	return session, nil
}

func (d mongoDBRepo) BeginTx(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := d.NewSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)
	return mongo.WithSession(ctx, sess, func(sessionContext mongo.SessionContext) error {
		if err := fn(sessionContext); err != nil {
			return err
		}
		return sess.CommitTransaction(sessionContext)
	})
}

func (b mongoColRepo) NewSession() (mongo.Session, error) {
	session, err := b.collection.Database().Client().StartSession()
	if err != nil {
		return nil, err
	}
	if err := session.StartTransaction(); err != nil {
		return nil, err
	}
	return session, nil
}

func (b mongoColRepo) BeginTx(ctx context.Context, fn func(ctx context.Context) error) error {
	sess, err := b.NewSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)
	return mongo.WithSession(ctx, sess, func(sessionContext mongo.SessionContext) error {
		if err := fn(sessionContext); err != nil {
			return err
		}
		return sess.CommitTransaction(sessionContext)
	})
}
