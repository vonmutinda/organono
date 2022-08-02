package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/vonmutinda/organono/app/logger"
	"github.com/vonmutinda/organono/app/utils"
	"github.com/vonmutinda/organono/app/web/auth"
	"github.com/vonmutinda/organono/app/web/ctxhelper"
)

const requestIDHeaderKey = "x-request-id"
const userAgentHeaderKey = "user-agent"

func DefaultMiddlewares(
	sessionAuthenticator auth.SessionAuthenticator,
) []gin.HandlerFunc {

	return []gin.HandlerFunc{
		securer(),
		corsMiddleware(),

		setRequestId(),
		logHTTPRequest(),
		setupAppContext(sessionAuthenticator),

		panicRecovery(),
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-ORGANONO-Token, X-SET-ORGANONO-ID, X-ORGANONO-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token, Authorization, X-Requested-With, X-ORGANONO-Token, X-SET-ORGANONO-ID, X-ORGANONO-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func panicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {

			if err := recover(); err != nil {

				requestID := ctxhelper.RequestId(c.Request.Context())

				c.Writer.WriteHeader(http.StatusInternalServerError)
				logger.Errorf("Failed to recover from panic: %v", err)
				debug.PrintStack()

				if os.Getenv("ENVIRONMENT") != "production" {
					fmt.Fprintf(
						c.Writer,
						`{"error":"panic: %s","details":"See logs for more information (%s)."}`,
						err,
						requestID,
					)
				} else {

					fmt.Fprintf(
						c.Writer,
						`{"error":"internal server error (%s)"}`,
						ctxhelper.RequestId(c.Request.Context()),
					)
				}
			}
		}()

		c.Next()
	}
}

func logHTTPRequest() gin.HandlerFunc {
	return func(c *gin.Context) {

		logger.Info(
			fmt.Sprintf("[%v] %s %s%s %s",
				ctxhelper.RequestId(c.Request.Context()),
				c.Request.Method,
				c.Request.Host,
				c.Request.RequestURI,
				c.Request.Proto,
			),
		)

		c.Next()
	}
}

func securer() gin.HandlerFunc {
	return secure.Secure(secure.Options{
		SSLRedirect:          strings.ToLower(os.Getenv("FORCE_SSL")) == "true",
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	})
}

func setRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := utils.GenerateUUID()
		ctx := ctxhelper.WithRequestId(c.Request.Context(), requestId)
		c.Request = c.Request.WithContext(ctx)
		c.Header(requestIDHeaderKey, requestId)
		c.Next()
	}
}

func setupAppContext(
	sessionAuthenticator auth.SessionAuthenticator,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		ipAddress, err := getIP(c.Request)
		if err != nil {
			logger.Warnf("Unable to parse ipAddress %v from remote address: %v", c.Request.RemoteAddr, err)
		} else {
			ctx = ctxhelper.WithIpAddress(ctx, ipAddress)
		}

		userAgent := c.Request.Header.Get(userAgentHeaderKey)
		ctx = ctxhelper.WithUserAgent(ctx, userAgent)

		tokenInfo, err := sessionAuthenticator.TokenInfoFromRequest(c.Request)
		if err != nil {

			wrappedError := utils.NewError(
				err,
				"Failed to parse token info from request",
			).Notify()

			if wrappedError.Err() != auth.ErrTokenNotProvided {
				wrappedError.LogErrorMessages()
			}
		} else {
			ctx = ctxhelper.WithTokenInfo(ctx, tokenInfo)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-FORWARDED-FOR")
	forwarded := strings.Split(ips, ",")

	if ips != "" && len(forwarded) > 0 {
		return forwarded[0], nil
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}
