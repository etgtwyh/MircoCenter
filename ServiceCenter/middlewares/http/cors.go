package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func DefaultCors() gin.HandlerFunc {
	return cors.Default()
}
