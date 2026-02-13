package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"markoni23/url-shortener/internal/domain"
	"markoni23/url-shortener/internal/repository"
	"math/rand"
	"time"
)

type LinkService struct {
	basePath   string
	repository repository.LinkRepository
}

func NewLinkService(basePath string, repository repository.LinkRepository) *LinkService {
	return &LinkService{
		basePath:   basePath,
		repository: repository,
	}
}

const SHORT_NAME_LENGTH = 8

func (l *LinkService) GetLinksList(ctx context.Context) ([]domain.Link, error) {
	return l.repository.FindAll(ctx)
}

func (l *LinkService) GetLink(ctx context.Context, id int64) (domain.Link, error) {
	link, err := l.repository.Get(ctx, id)
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

func (l *LinkService) UpdateLink(ctx context.Context, id int64, originalUrl, shortName string) (domain.Link, error) {
	link, err := l.repository.Update(ctx, id, repository.UpdateLinkDTO{
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
	return link, nil
}

func (l *LinkService) DeleteLink(ctx context.Context, id int64) error {
	return l.repository.Delete(ctx, id)
}

func (l *LinkService) CreateLink(ctx context.Context, originalUrl, shortName string) (domain.Link, error) {
	if shortName == "" {
		shortName = GenerateShortName()
	}
	shortLink := fmt.Sprintf("%s/%s", l.basePath, GenerateShortName())

	return l.repository.Create(ctx, repository.CreateLinkDTO{
		OriginalUrl: &originalUrl,
		ShortName:   &shortName,
		ShortUrl:    &shortLink,
	})
}

func GenerateShortName() string {
	alphabet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]byte, SHORT_NAME_LENGTH)
	for i := range SHORT_NAME_LENGTH {
		res[i] = alphabet[rand.Int()%len(alphabet)]
	}
	return string(res)
}
