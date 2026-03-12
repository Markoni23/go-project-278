package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"

	"markoni23/url-shortener/internal/model"
	"markoni23/url-shortener/internal/sqlcdb"
)

type service struct {
	basePath string
	queries  *sqlcdb.Queries
}

func NewService(basePath string, queries *sqlcdb.Queries) *service {
	return &service{
		basePath: basePath,
		queries:  queries,
	}
}

func (s *service) Count(ctx context.Context) (int64, error) {
	return s.queries.GetLinksCount(ctx)
}

func (s *service) GetAll(ctx context.Context, from, to int64) ([]model.Link, error) {
	if from < 0 || to <= 0 {
		return []model.Link{}, errors.New("from and to must be greater than zero")
	}

	if from >= to {
		return []model.Link{}, errors.New("from must be less than to")
	}

	limit := to - from + 1
	offset := from
	linksRaw, err := s.queries.GetLinks(ctx, sqlcdb.GetLinksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		return []model.Link{}, err
	}

	res := make([]model.Link, len(linksRaw))
	for i, raw := range linksRaw {
		res[i] = s.rawToModel(raw)
	}

	return res, nil
}

func (s *service) Get(ctx context.Context, id int64) (model.Link, error) {
	link, err := s.queries.GetLink(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.Link{}, &model.LinkNotFoundError{}
		default:
			return model.Link{}, nil
		}
	}
	return s.rawToModel(link), nil
}

func (s *service) GetLinkByShortName(ctx context.Context, shortName string) (model.Link, error) {
	link, err := s.queries.GetLinkByShortName(ctx, sql.NullString{String: shortName, Valid: true})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.Link{}, &model.LinkNotFoundError{}
		default:
			return model.Link{}, nil
		}
	}
	return s.rawToModel(link), nil
}

func (s *service) Update(ctx context.Context, id int64, originalURL, shortName string) (model.Link, error) {
	res, err := s.queries.UpdateLink(ctx, sqlcdb.UpdateLinkParams{
		ID:          id,
		OriginalUrl: sql.NullString{String: originalURL, Valid: true},
		ShortName:   sql.NullString{String: shortName, Valid: true},
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.Link{}, &model.LinkNotFoundError{}
		default:
			return model.Link{}, nil
		}
	}
	return s.rawToModel(res), nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	_, err := s.queries.GetLink(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &model.LinkNotFoundError{}
		default:
			return err
		}
	}
	return s.queries.DeleteLink(ctx, id)
}

func (s *service) Create(ctx context.Context, originalURL, shortName string) (model.Link, error) {
	if shortName == "" {
		shortName = GenerateShortName()
	}

	res, err := s.queries.CreateLink(ctx, sqlcdb.CreateLinkParams{
		OriginalUrl: sql.NullString{String: originalURL, Valid: true},
		ShortName:   sql.NullString{String: shortName, Valid: true},
	})

	if err != nil {
		return model.Link{}, err
	}

	return s.rawToModel(res), nil
}

const ShortNameLength = 8

func GenerateShortName() string {
	alphabet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	res := make([]byte, ShortNameLength)
	for i := range ShortNameLength {
		res[i] = alphabet[rand.IntN(len(alphabet))]
	}
	return string(res)
}

func (s *service) rawToModel(raw sqlcdb.Link) model.Link {
	shortUrl := fmt.Sprintf("%s/r/%s", s.basePath, raw.ShortName.String)
	return model.Link{
		ID:          raw.ID,
		OriginalUrl: raw.OriginalUrl.String,
		ShortName:   raw.ShortName.String,
		ShortUrl:    shortUrl,
	}
}
