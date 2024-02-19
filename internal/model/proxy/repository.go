package proxy

import (
	"context"
	"database/sql"
	"log"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetActiveProxy(ctx context.Context) (*Proxy, error) {
	query := `
			SELECT * FROM proxies
			WHERE isactive = true
			ORDER BY id
			LIMIT 1;
		`

	proxy := &Proxy{}

	row := r.db.QueryRowContext(ctx, query)
	err := row.Scan(&proxy.Id, &proxy.Body, &proxy.Active)
	if err != nil {
		log.Fatal("Error scanning row:", err)
		return nil, err
	}
	return proxy, err
}

func (r *repo) UpdateProxy(ctx context.Context, body string) error {
	query := `
			UPDATE proxies 
			SET isactive = true
			WHERE body = $1;
		`

	_, err := r.db.ExecContext(ctx, query, body)
	if err != nil {
		log.Fatal("Error scanning row:", err)
	}
	return nil
}
