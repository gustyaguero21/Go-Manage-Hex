package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UrlMapping(s *gin.Engine) {
	api := s.Group("/api/go-manage-hex")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
