package linkvisit

import (
	"fmt"
	"markoni23/url-shortener/internal/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type VisitService interface {
	GetAll(ctx context.Context, from, to int64) ([]model.LinkVisit, error)
	Visit(ctx *gin.Context, link model.Link) error
	Count(ctx context.Context) (int64, error)
}

type LinkService interface {
	GetLinkByShortName(ctx context.Context, shortName string) (model.Link, error)
}

type handler struct {
	visitService VisitService
	linkService  LinkService
}

func NewHandler(visitService VisitService, linkService LinkService) *handler {
	return &handler{
		visitService: visitService,
		linkService:  linkService,
	}
}

func (h *handler) VisistLink(ctx *gin.Context) {
	code := ctx.Param("code")

	link, err := h.linkService.GetLinkByShortName(ctx, code)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.visitService.Visit(ctx, link); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *handler) GetVisits(ctx *gin.Context) {
	rangeString := ctx.DefaultQuery("range", "[0,10]")

	rangeWithoutBrackets := strings.Trim(rangeString, "[]")
	fromToSlice := strings.Split(rangeWithoutBrackets, ",")

	if len(fromToSlice) != 2 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong range format"})
	}

	from, err := strconv.Atoi(strings.TrimSpace(fromToSlice[0]))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' value"})
	}

	to, err := strconv.Atoi(strings.TrimSpace(fromToSlice[1]))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' value"})
	}

	res, err := h.visitService.GetAll(ctx, int64(from), int64(to))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	count, err := h.visitService.Count(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	headerValue := fmt.Sprintf("links_visits %d-%d/%d", from, to, count)
	ctx.Header("Content-Range", headerValue)

	ctx.JSON(http.StatusOK, res)
}
