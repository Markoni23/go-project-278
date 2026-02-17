package handler

import (
	"github.com/gin-gonic/gin"
)

type LinkHandler interface {
	GetLinksList(ctx *gin.Context)
	GetLink(ctx *gin.Context)
	CreateLink(ctx *gin.Context)
	UpdateLink(ctx *gin.Context)
	DeleteLink(ctx *gin.Context)
}

func RegisterRoutes(parentGroup *gin.RouterGroup, handler LinkHandler) {
	linksRoutes := parentGroup.Group("/links")
	{
		linksRoutes.GET("/", handler.GetLinksList)
		linksRoutes.POST("/", handler.CreateLink)
		linksRoutes.GET("/:id", handler.GetLink)
		linksRoutes.PUT("/:id", handler.UpdateLink)
		linksRoutes.DELETE("/:id", handler.DeleteLink)
	}
}
