package server

import (
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/infrastructure/auth"
	"go-manage-hex/internal/infrastructure/db"
	"time"

	"log"

	service "go-manage-hex/internal/app/user"
	repository "go-manage-hex/internal/infrastructure/db/user"
	handler "go-manage-hex/internal/infrastructure/http/handler/user"
	middleware "go-manage-hex/internal/infrastructure/http/middleware"

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

	authService := auth.NewJWTService(config.GetJwtSecret(), time.Hour*1)
	middleware := middleware.NewMiddleware(authService)

	userHandler := handler.NewUserHandler(userService, authService)

	api := s.Group(config.BaseURL)

	api.POST("/create", userHandler.CreateUserHandler)
	api.POST("/login", userHandler.LoginUser)

	protected := api.Group("/")
	protected.Use(middleware.RequireAuth)

	protected.GET("/search", userHandler.SearchUserHandler)

	protected.DELETE("/delete", userHandler.DeleteUserHandler)

	protected.PATCH("/update", userHandler.UpdateUserHandler)

	protected.PATCH("/change-password", userHandler.ChangePwdHandler)

}
