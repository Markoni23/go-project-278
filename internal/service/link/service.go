package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"markoni23/url-shortener/internal/domain"
	"markoni23/url-shortener/internal/dto"
	"math/rand"
	"time"
)

type LinkRepository interface {
	FindAll(ctx context.Context) ([]domain.Link, error)
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

func (s *service) GetAll(ctx context.Context) ([]domain.Link, error) {
	return s.repository.FindAll(ctx)
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

func (s *service) Update(ctx context.Context, id int64, originalUrl, shortName string) (domain.Link, error) {
	res, err := s.repository.Update(ctx, id, dto.UpdateLinkDTO{
		OriginalUrl: &originalUrl,
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

func (s *service) Create(ctx context.Context, originalUrl, shortName string) (domain.Link, error) {
	if shortName == "" {
		shortName = GenerateShortName()
	}
	shortLink := fmt.Sprintf("%s/%s", s.basePath, GenerateShortName())

	return s.repository.Create(ctx, dto.CreateLinkDTO{
		OriginalUrl: &originalUrl,
		ShortName:   &shortName,
		ShortUrl:    &shortLink,
	})
}

const SHORT_NAME_LENGTH = 8

func GenerateShortName() string {
	alphabet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]byte, SHORT_NAME_LENGTH)
	for i := range SHORT_NAME_LENGTH {
		res[i] = alphabet[rand.Int()%len(alphabet)]
	}
	return string(res)
}
