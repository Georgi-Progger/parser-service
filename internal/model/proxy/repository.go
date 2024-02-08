package proxy

import (
	"context"
	"database/sql"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func(r *repo) GetActiveProxy(ctx context.Context){
	
}
