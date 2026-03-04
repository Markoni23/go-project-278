package linkvisit

import (
	"context"
	"database/sql"
	"errors"
	"markoni23/url-shortener/internal/model"
	"markoni23/url-shortener/internal/sqlcdb"
	"net/http"

	"github.com/gin-gonic/gin"
)

type service struct {
	queries *sqlcdb.Queries
}

func NewService(queries *sqlcdb.Queries) *service {
	return &service{
		queries: queries,
	}
}

func (s *service) Count(ctx context.Context) (int64, error) {
	return s.queries.CountLinkVisits(ctx)
}

func (s *service) GetAll(ctx context.Context, from, to int64) ([]model.LinkVisit, error) {
	if from < 0 || to <= 0 {
		return []model.LinkVisit{}, errors.New("from and to must be greater than zero")
	}

	if from >= to {
		return []model.LinkVisit{}, errors.New("from must be less than to")
	}

	limit := to - from + 1
	offset := from
	visits, err := s.queries.GetAllLinkVisits(ctx, sqlcdb.GetAllLinkVisitsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return []model.LinkVisit{}, err
	}

	res := make([]model.LinkVisit, len(visits))

	for i, raw := range visits {
		res[i] = s.rawToModel(raw)
	}
	return res, nil
}

func (s *service) Visit(ctx *gin.Context, link model.Link) error {
	userAgent := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")
	params := sqlcdb.CreateLinkVisitParams{
		LinkID:    link.ID,
		Ip:        ctx.ClientIP(),
		UserAgent: sql.NullString{String: userAgent, Valid: true},
		Referer:   sql.NullString{String: referer, Valid: true},
		Status:    http.StatusFound,
	}

	_, err := s.queries.CreateLinkVisit(ctx, params)
	if err != nil {
		return err
	}

	ctx.Redirect(http.StatusFound, link.OriginalUrl)

	return nil
}

func (s *service) rawToModel(raw sqlcdb.LinkVisit) model.LinkVisit {
	return model.LinkVisit{
		ID:        raw.ID,
		LinkId:    raw.LinkID,
		Ip:        raw.Ip,
		UserAgent: &raw.UserAgent.String,
		Referer:   &raw.Referer.String,
		Status:    int64(raw.Status),
		CreatedAt: raw.CreatedAt.Time,
	}
}
