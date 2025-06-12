package server

import (
	"context"
	"database/sql"
	"go-manage-hex/cmd/config"
	"go-manage-hex/internal/infrastructure/db"

	"log"
	"net/http"

	userservice "go-manage-hex/internal/app/user"
	entity "go-manage-hex/internal/core/user"
	repository "go-manage-hex/internal/infrastructure/db/user"

	"github.com/gin-gonic/gin"
)

func UrlMapping(s *gin.Engine) {

	db, err := db.DatabaseConn()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserMysql(db)
	repo.CreateTable(config.GetMysqlTable())

	//testSearch(db)
	//testCreate(db)
	//testDelete(db)
	//testUpdate(db)
	//testChangePwd(db)

	api := s.Group(config.BaseURL)

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func testSearch(db *sql.DB) {

	repo := repository.NewUserMysql(db)

	service := userservice.NewUserService(repo)

	search, searchErr := service.SearchUser(context.Background(), "johndoe")
	if searchErr != nil {
		log.Fatal(searchErr)
	}
	log.Print(search)
}

func testCreate(db *sql.DB) {
	repo := repository.NewUserMysql(db)

	service := userservice.NewUserService(repo)

	user := entity.User{
		Name:     "John",
		LastName: "Doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
		Password: "Password1234567",
	}

	created, createdErr := service.CreateUser(context.Background(), user)
	if createdErr != nil {
		log.Fatal(createdErr)
	}

	log.Print(created)
}

func testDelete(db *sql.DB) {
	repo := repository.NewUserMysql(db)

	service := userservice.NewUserService(repo)

	username := "johndoe"

	if deleteErr := service.DeleteUser(context.Background(), username); deleteErr != nil {
		log.Fatal(deleteErr)
	}
	log.Print("user deleted successfully")
}

func testUpdate(db *sql.DB) {
	repo := repository.NewUserMysql(db)

	service := userservice.NewUserService(repo)

	username := "johndoe"

	update := entity.User{
		Name:     "Johncito",
		LastName: "Doecito",
		Email:    "johncitodoecito@example.com",
	}

	updated, updateErr := service.UpdateUser(context.Background(), username, update)
	if updateErr != nil {
		log.Fatal(updateErr)
	}
	log.Print(updated)
}

func testChangePwd(db *sql.DB) {
	repo := repository.NewUserMysql(db)

	service := userservice.NewUserService(repo)

	username := "johndoe"
	newPwd := "NewPassword1234"

	if err := service.ChangeUserPwd(context.Background(), newPwd, username); err != nil {
		log.Fatal(err)
	}
	log.Print("password changed successfully")
}
