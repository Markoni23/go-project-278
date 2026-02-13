package repository

import (
	"context"
	"markoni23/url-shortener/internal/domain"
)

type LinkRepository interface {
	FindAll(ctx context.Context) ([]domain.Link, error)
	Get(ctx context.Context, id int64) (domain.Link, error)
	Create(ctx context.Context, dto CreateLinkDTO) (domain.Link, error)
	Update(ctx context.Context, id int64, dto UpdateLinkDTO) (domain.Link, error)
	Delete(ctx context.Context, id int64) error
}

type CreateLinkDTO struct {
	OriginalUrl *string
	ShortName   *string
	ShortUrl    *string
}

type UpdateLinkDTO struct {
	ShortName   *string
	OriginalUrl *string
}
