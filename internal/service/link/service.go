package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"

	"markoni23/url-shortener/internal/domain"
	"markoni23/url-shortener/internal/dto"
)

type LinkRepository interface {
	Count(ctx context.Context) (int64, error)
	FindAll(ctx context.Context, dto dto.GetLinksDTO) ([]domain.Link, error)
	Get(ctx context.Context, id int64) (domain.Link, error)
	Create(ctx context.Context, dto dto.CreateLinkDTO) (domain.Link, error)
	Update(ctx context.Context, id int64, dto dto.UpdateLinkDTO) (domain.Link, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	basePath   string
	repository LinkRepository
}

func NewService(basePath string, repository LinkRepository) *service {
	return &service{
		basePath:   basePath,
		repository: repository,
	}
}

func (s *service) Count(ctx context.Context) (int64, error) {
	return s.repository.Count(ctx)
}

func (s *service) GetAll(ctx context.Context, from, to int64) ([]domain.Link, error) {
	if from <= 0 || to <= 0 {
		return []domain.Link{}, errors.New("from and to must be greater than zero")
	}

	if from >= to {
		return []domain.Link{}, errors.New("from must be less than to")
	}

	limit := to - from + 1
	offset := from - 1
	return s.repository.FindAll(ctx, dto.GetLinksDTO{
		Limit:  &limit,
		Offset: &offset,
	})
}

func (s *service) Get(ctx context.Context, id int64) (domain.Link, error) {
	link, err := s.repository.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.Link{}, &domain.LinkNotFoundError{}
		default:
			return domain.Link{}, nil
		}
	}
	return link, nil
}

func (s *service) Update(ctx context.Context, id int64, originalURL, shortName string) (domain.Link, error) {
	res, err := s.repository.Update(ctx, id, dto.UpdateLinkDTO{
		OriginalUrl: &originalURL,
		ShortName:   &shortName,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.Link{}, &domain.LinkNotFoundError{}
		default:
			return domain.Link{}, nil
		}
	}
	return res, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) Create(ctx context.Context, originalURL, shortName string) (domain.Link, error) {
	if shortName == "" {
		shortName = GenerateShortName()
	}
	shortLink := fmt.Sprintf("%s/%s", s.basePath, GenerateShortName())

	return s.repository.Create(ctx, dto.CreateLinkDTO{
		OriginalUrl: &originalURL,
		ShortName:   &shortName,
		ShortUrl:    &shortLink,
	})
}

const ShortNameLength = 8

func GenerateShortName() string {
	alphabet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	random := rand.New(&rand.Rand{})
	res := make([]byte, ShortNameLength)
	for i := range ShortNameLength {
		res[i] = alphabet[random.Int()%len(alphabet)]
	}
	return string(res)
}
