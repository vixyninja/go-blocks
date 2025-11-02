package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey string

const (
	RequestIDKey    ContextKey = "requestID"
	RequestIDHeader string     = "X-Request-ID"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(string(RequestIDKey), requestID)

		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(string(RequestIDKey)); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

func GetRequestIDFromContext(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
