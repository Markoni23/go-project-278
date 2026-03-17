package link

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"markoni23/url-shortener/internal/model"
	"markoni23/url-shortener/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Count(ctx context.Context) (int64, error)
	GetAll(ctx context.Context, from, to int64) ([]model.Link, error)
	Get(ctx context.Context, id int64) (model.Link, error)
	Create(ctx context.Context, originalURL, shortName string) (model.Link, error)
	Update(ctx context.Context, id int64, originalURL, shortName string) (model.Link, error)
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
	rangeString := ctx.DefaultQuery("range", "[0,10]")

	rangeWithoutBrackets := strings.Trim(rangeString, "[]")
	fromToSlice := strings.Split(rangeWithoutBrackets, ",")

	if len(fromToSlice) != 2 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong range format"})
		return
	}

	from, err := strconv.Atoi(strings.TrimSpace(fromToSlice[0]))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' value"})
		return
	}

	to, err := strconv.Atoi(strings.TrimSpace(fromToSlice[1]))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' value"})
		return
	}

	res, err := l.service.GetAll(ctx, int64(from), int64(to))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := l.service.Count(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	headerValue := fmt.Sprintf("links %d-%d/%d", from, to, count)
	ctx.Header("Content-Range", headerValue)

	ctx.JSON(http.StatusOK, res)
}

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required,url"`
	ShortName   string `json:"short_name" binding:"omitempty,min=3,max=32"`
}

func (h *handler) CreateLink(ctx *gin.Context) {
	var r CreateLinkRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			ctx.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{
				Errors: utils.FormatValidationErrors(err),
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, utils.SimpleErrorResponse{
			Error: "invalid request",
		})
		return
	}

	link, err := h.service.Create(ctx, r.OriginalUrl, r.ShortName)
	if err != nil {
		if utils.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{
				Errors: utils.FormatDuplicateKeyError(err, "short_name"),
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create link"})
		return
	}

	ctx.JSON(http.StatusCreated, link)
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
		if errors.Is(err, &model.LinkNotFoundError{}) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, link)
}

type UpdateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required,url"`
	ShortName   string `json:"short_name" binding:"omitempty,min=3,max=32"`
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
		if _, ok := err.(validator.ValidationErrors); ok {
			ctx.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{
				Errors: utils.FormatValidationErrors(err),
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, utils.SimpleErrorResponse{
			Error: "invalid request",
		})
		return
	}

	link, err := h.service.Update(ctx, id, req.OriginalUrl, req.ShortName)
	if err != nil {
		if errors.Is(err, &model.LinkNotFoundError{}) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}

		if utils.IsDuplicateKeyError(err) {
			ctx.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{
				Errors: utils.FormatDuplicateKeyError(err, "short_name"),
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update link"})
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
		if errors.Is(err, &model.LinkNotFoundError{}) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
