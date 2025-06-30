package server

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/infrastructure/db"

	"log"
	"net/http"

	service "go-manage-hex/internal/app/user"
	repository "go-manage-hex/internal/infrastructure/db/user"
	handler "go-manage-hex/internal/infrastructure/http/handler/user"

	"github.com/gin-gonic/gin"
)

func UrlMapping(s *gin.Engine) {

	db, err := db.DatabaseConn()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserMysql(db)
	userRepo.CreateTable(config.GetMysqlTable())

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	api := s.Group(config.BaseURL)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.GET("/search", userHandler.SearchUserHandler)

	api.POST("/create", userHandler.CreateUserHandler)

	api.DELETE("/delete", userHandler.DeleteUserHandler)

	api.PATCH("/update", userHandler.UpdateUserHandler)

	api.PATCH("/change-password", userHandler.ChangePwdHandler)

}
