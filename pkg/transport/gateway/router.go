package gateway

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/scrynaxx/tr-core/pkg/tracing"
	"github.com/scrynaxx/tr-core/pkg/transport"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc/codes"
)

type revokeRepository interface {
	IsRevoked(ctx context.Context, sessionID uuid.UUID) (bool, error)
}

func newRouter(
	mux *runtime.ServeMux,
	revokeRepo revokeRepository,
	verifier *tokenVerifier,
	publicRoutes []PublicRoute,
	middleware ...gin.HandlerFunc,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	router := gin.New()
	router.Use(recoveryMiddleware)
	router.Use(middleware...)
	router.Use(otelgin.Middleware("gateway"))
	router.Use(requestIDMiddleware)
	router.Use(authMiddleware(revokeRepo, verifier, newRouteMatcher(publicRoutes)))
	router.Any("/*path", gin.WrapH(mux))

	return router
}

func recoveryMiddleware(c *gin.Context) {
	defer func() {
		if recovered := recover(); recovered != nil {
			slog.Error("[gateway] http router panic recovered",
				slog.Any("panic", recovered),
				slog.String("path", c.Request.URL.Path),
				slog.String("method", c.Request.Method),
				slog.String("remote_addr", c.Request.RemoteAddr),
				slog.String("user_agent", c.Request.UserAgent()),
			)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    codes.Internal,
				"message": "internal server error",
			})
		}
	}()

	c.Next()
}

func requestIDMiddleware(c *gin.Context) {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		c.Next()
		return
	}

	ctx := tracing.WithRequestID(c.Request.Context(), requestID)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func authMiddleware(revokeRepo revokeRepository, verifier *tokenVerifier, publicRoutes routeMatcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		if publicRoutes.Match(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		identity, err := verifier.Validate(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		revoked, err := revokeRepo.IsRevoked(c.Request.Context(), identity.SessionID)
		if err != nil {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}
		if revoked {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Request.Header.Set(transport.EmployeeIDHeader, identity.EmployeeID.String())
		c.Request.Header.Set(transport.SessionIDHeader, identity.SessionID.String())

		c.Next()
	}
}
