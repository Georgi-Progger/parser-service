package repositories

import (
	"context"
	"database/sql"
	"log"

	"main.go/internal/model/proxy"
)

type proxyRepository struct {
	db *sql.DB
}

func NewProxyRepository(db *sql.DB) *proxyRepository {
	return &proxyRepository{
		db: db,
	}
}
func (r *proxyRepository) GetActiveProxy(ctx context.Context) (*proxy.Proxy, error) {
	query := `
			SELECT * FROM proxies
			WHERE isactive = true
			ORDER BY id
			LIMIT 1;
		`

	proxy := &proxy.Proxy{}

	row := r.db.QueryRowContext(ctx, query)
	err := row.Scan(&proxy.Id, &proxy.Body, &proxy.Active)
	if err != nil {
		log.Fatal("Error scanning row:", err)
		return nil, err
	}
	return proxy, err
}

func (r *proxyRepository) BlockProxy(ctx context.Context, body string) error {
	query := `
			UPDATE proxies 
			SET isactive = false
			WHERE body = $1;
		`

	_, err := r.db.ExecContext(ctx, query, body)
	if err != nil {
		log.Fatal("Error scanning row:", err)
	}
	return nil
}
