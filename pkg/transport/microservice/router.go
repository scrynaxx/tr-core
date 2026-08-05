package microservice

import (
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	router := gin.New()
	router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		slog.Error("[service] http router panic recovered",
			slog.String("path", c.Request.URL.Path),
			slog.Any("error", err),
			slog.String("stack", string(debug.Stack())),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	return router
}
