package app

import (
	"fmt"
	"log"

	"markoni23/url-shortener/internal/config"
	"markoni23/url-shortener/internal/db"
	"markoni23/url-shortener/internal/handler"
	linkHandler "markoni23/url-shortener/internal/handler/link"
	linkRepository "markoni23/url-shortener/internal/repository/link"
	linkService "markoni23/url-shortener/internal/service/link"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func Run(cfg config.Config) error {
	db, err := db.InitDB(cfg.Database.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))

	router.GET("/", func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
				hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			})
		}
		ctx.Status(http.StatusOK)
	})

	linkRepo := linkRepository.NewDBLinkRepository(db)
	linkSvc := linkService.NewService(cfg.Server.BasePath, linkRepo)
	linkHand := linkHandler.NewHandler(linkSvc)

	handler.RegisterRoutes(&router.RouterGroup, linkHand)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router.Run(":" + cfg.Server.Port)
}
