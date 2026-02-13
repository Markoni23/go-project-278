package db

import (
	"context"
	"database/sql"
	"markoni23/url-shortener/internal/domain"
	"markoni23/url-shortener/internal/repository"
	"markoni23/url-shortener/internal/sqlcdb"
)

type DBLinkRepositry struct {
	queries *sqlcdb.Queries
}

func NewDBLinkRepository(pool *sql.DB) *DBLinkRepositry {
	return &DBLinkRepositry{
		queries: sqlcdb.New(pool),
	}
}

func (d *DBLinkRepositry) FindAll(ctx context.Context) ([]domain.Link, error) {
	links, err := d.queries.GetLinks(ctx)
	if err != nil {
		return []domain.Link{}, err
	}
	res := make([]domain.Link, len(links))
	for i, raw := range links {
		res[i] = sqlcdbToDomainLink(raw)
	}
	return res, nil
}

func (d *DBLinkRepositry) Get(ctx context.Context, id int64) (domain.Link, error) {
	link, err := d.queries.GetLink(ctx, id)
	if err != nil {
		return domain.Link{}, err
	}

	return sqlcdbToDomainLink(link), nil
}

func (d *DBLinkRepositry) Create(ctx context.Context, dto repository.CreateLinkDTO) (domain.Link, error) {
	link, err := d.queries.CreateLink(ctx, sqlcdb.CreateLinkParams{
		OriginalUrl: sql.NullString{String: *dto.OriginalUrl, Valid: dto.OriginalUrl != nil},
		ShortName:   sql.NullString{String: *dto.ShortName, Valid: dto.ShortName != nil},
		ShortUrl:    sql.NullString{String: *dto.ShortUrl, Valid: dto.ShortUrl != nil},
	})

	if err != nil {
		return domain.Link{}, err
	}

	return sqlcdbToDomainLink(link), nil
}

func (d *DBLinkRepositry) Update(ctx context.Context, id int64, dto repository.UpdateLinkDTO) (domain.Link, error) {
	params := sqlcdb.UpdateLinkParams{
		ID:          id,
		OriginalUrl: sql.NullString{String: *dto.OriginalUrl, Valid: dto.OriginalUrl != nil},
		ShortName:   sql.NullString{String: *dto.ShortName, Valid: dto.ShortName != nil},
	}
	link, err := d.queries.UpdateLink(ctx, params)
	if err != nil {
		return domain.Link{}, err
	}

	return sqlcdbToDomainLink(link), nil
}

func (d *DBLinkRepositry) Delete(ctx context.Context, id int64) error {
	return d.queries.DeleteLink(ctx, id)
}

func sqlcdbToDomainLink(raw sqlcdb.Link) domain.Link {
	return domain.Link{
		ID:          raw.ID,
		OriginalUrl: raw.OriginalUrl.String,
		ShortName:   raw.ShortName.String,
		ShortUrl:    raw.ShortUrl.String,
	}
}
