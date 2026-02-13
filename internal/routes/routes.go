package routes

import (
	"markoni23/url-shortener/internal/controller"
	"markoni23/url-shortener/internal/service"

	"github.com/gin-gonic/gin"
)

func LinkRoutes(parentGroup *gin.RouterGroup, service service.LinkService) {
	linksRoutes := parentGroup.Group("/links")
	linkController := controller.LinkController{Service: service}
	{
		linksRoutes.GET("/", linkController.GetLinksList)
		linksRoutes.POST("/", linkController.CreateLink)
		linksRoutes.GET("/:id", linkController.GetLink)
		linksRoutes.PUT("/:id", linkController.UpdateLink)
		linksRoutes.DELETE("/:id", linkController.DeleteLink)
	}
}
