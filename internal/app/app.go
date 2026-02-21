package app

import (
	"database/sql"
	"markoni23/url-shortener/internal/config"
	linkHandler "markoni23/url-shortener/internal/handler/link"
	linkRepository "markoni23/url-shortener/internal/repository/link"
	linkService "markoni23/url-shortener/internal/service/link"
	"net/http"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(cfg config.Config, db *sql.DB) error {
	router := gin.Default()

	if cfg.Server.SentryDSN != "" {
		router.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.Server.BasePath}
	if cfg.Server.FrontendUrl != "" {
		corsConfig.AllowOrigins = []string{cfg.Server.BasePath, cfg.Server.FrontendUrl}
	}

	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(corsConfig))

	linkRepo := linkRepository.NewDBLinkRepository(db)
	linkSvc := linkService.NewService(cfg.Server.BasePath, linkRepo)
	linkHand := linkHandler.NewHandler(linkSvc)

	apiGroup := router.Group("/api")
	{
		linksRoutes := apiGroup.Group("/links")
		{
			linksRoutes.GET("/", linkHand.GetLinksList)
			linksRoutes.POST("/", linkHand.CreateLink)
			linksRoutes.GET("/:id", linkHand.GetLink)
			linksRoutes.PUT("/:id", linkHand.UpdateLink)
			linksRoutes.DELETE("/:id", linkHand.DeleteLink)
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router.Run(":" + cfg.Server.Port)
}
