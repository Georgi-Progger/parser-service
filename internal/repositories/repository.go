package repositories

import (
	"context"
	"database/sql"

	annoucement "main.go/internal/model/annoucement"
	"main.go/internal/model/proxy"
)

type AnnoucementRepository interface {
	GetAnnoucements(ctx context.Context, page int) (*[]annoucement.Annoucement, error)
	SetAnnoucement(ctx context.Context, annoucementInfo annoucement.Annoucement) error
	LinkExists(ctx context.Context, link string) bool
}

type ProxyRepository interface {
	GetActiveProxy(ctx context.Context) (*proxy.Proxy, error)
	UpdateProxy(ctx context.Context, body string) error
}

type repository struct {
	AnnoucementRepository
	ProxyRepository
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		AnnoucementRepository: NewAnnoucementRepository(db),
		ProxyRepository:       NewProxyRepository(db),
	}
}
