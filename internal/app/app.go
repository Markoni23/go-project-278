package app

import (
	"context"
	"database/sql"
	"markoni23/url-shortener/internal/config"
	. "markoni23/url-shortener/internal/db"
	"markoni23/url-shortener/internal/routes"
	"markoni23/url-shortener/internal/service"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, cfg config.Config, db *sql.DB) error {
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

	linkRepo := NewDBLinkRepository(db)
	linkService := service.NewLinkService(cfg.Server.BasePath, linkRepo)
	routes.LinkRoutes(&router.RouterGroup, *linkService)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.Run(":" + cfg.Server.Port)

	return nil
}
