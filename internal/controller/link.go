package controller

import (
	"errors"
	"markoni23/url-shortener/internal/domain"
	"markoni23/url-shortener/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	Service service.LinkService
}

func NewLinkController(service service.LinkService) *LinkController {
	return &LinkController{
		Service: service,
	}
}

func (l *LinkController) GetLinksList(ctx *gin.Context) {

	res, err := l.Service.GetLinksList(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, res)
}

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortName   string `json:"short_name,omitempty"`
}

func (l *LinkController) CreateLink(ctx *gin.Context) {
	var r CreateLinkRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := l.Service.CreateLink(ctx, r.OriginalUrl, r.ShortName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (l *LinkController) GetLink(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := l.Service.GetLink(ctx, id)
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

func (l *LinkController) UpdateLink(ctx *gin.Context) {
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

	link, err := l.Service.UpdateLink(ctx, id, req.OriginalUrl, req.ShortName)
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

func (l *LinkController) DeleteLink(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := l.Service.DeleteLink(ctx, id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
