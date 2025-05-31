package server

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/infrastructure/db"
	"log"
	"net/http"

	//user "go-manage-hex/internal/infrastructure/db/user"

	"github.com/gin-gonic/gin"
)

func UrlMapping(s *gin.Engine) {

	_, err := db.DatabaseConn()
	if err != nil {
		log.Fatal(err)
	}

	// repo := user.NewUserMysql(db)
	// if tableErr := repo.CreateTable(config.GetMysqlTable()); tableErr != nil {
	// 	log.Fatal(tableErr)
	// }

	api := s.Group(config.BaseURL)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
