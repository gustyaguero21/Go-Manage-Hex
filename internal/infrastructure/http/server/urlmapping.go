package server

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/infrastructure/db"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UrlMapping(s *gin.Engine) {

	_, err := db.DatabaseConn()
	if err != nil {
		log.Fatal(err)
	}

	api := s.Group(config.BaseURL)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
