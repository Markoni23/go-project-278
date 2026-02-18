package link

import (
	"context"
	"errors"
	"fmt"
	"markoni23/url-shortener/internal/domain"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Count(ctx context.Context) (int64, error)
	GetAll(ctx context.Context, from, to int64) ([]domain.Link, error)
	Get(ctx context.Context, id int64) (domain.Link, error)
	Create(ctx context.Context, originalUrl, shortName string) (domain.Link, error)
	Update(ctx context.Context, id int64, originalUrl, shortName string) (domain.Link, error)
	Delete(ctx context.Context, id int64) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (l *handler) GetLinksList(ctx *gin.Context) {
	rangeString := ctx.DefaultQuery("range", "[1,10]")

	rangeWithoutBrackets := strings.Trim(rangeString, "[]")
	fromToSlice := strings.Split(rangeWithoutBrackets, ",")

	if len(fromToSlice) != 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong range format"})
		return
	}

	from, err := strconv.Atoi(strings.TrimSpace(fromToSlice[0]))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid 'from' value"})
		return
	}

	to, err := strconv.Atoi(strings.TrimSpace(fromToSlice[1]))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid 'to' value"})
		return
	}

	res, err := l.service.GetAll(ctx, int64(from), int64(to))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	count, err := l.service.Count(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	
	headerValue := fmt.Sprintf("links %d-%d/%d", from, to, count)
	ctx.Header("Content-Range", headerValue)

	ctx.JSON(http.StatusOK, res)
}

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name,omitempty"`
}

func (h *handler) CreateLink(ctx *gin.Context) {
	var r CreateLinkRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.service.Create(ctx, r.OriginalUrl, r.ShortName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *handler) GetLink(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := h.service.Get(ctx, id)
	if err != nil {
		if errors.Is(err, &domain.LinkNotFoundError{}) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, link)
}

type UpdateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name" binding:"required"`
}

func (h *handler) UpdateLink(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateLinkRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	link, err := h.service.Update(ctx, id, req.OriginalUrl, req.ShortName)
	if err != nil {
		if errors.Is(err, &domain.LinkNotFoundError{}) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, link)
}

func (h *handler) DeleteLink(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
