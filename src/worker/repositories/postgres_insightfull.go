package repositories

import (
	"context"
	"database/sql"
)

type PostgresInsightfullRepository interface {
	FindStoresByMerchantID(ctx context.Context, merchantID int64) error
}

type postgresInsightfullRepoImpl struct {
	orm *sql.DB
}

func NewPostgresInsightfullRepository(orm *sql.DB) PostgresInsightfullRepository {
	return &postgresInsightfullRepoImpl{
		orm: orm,
	}
}

func (r *postgresInsightfullRepoImpl) FindStoresByMerchantID(ctx context.Context, merchantID int64) error {

	return nil
}
