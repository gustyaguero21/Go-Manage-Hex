package server

import "github.com/gin-gonic/gin"

func StartServer() *gin.Engine {
	serve := gin.Default()

	UrlMapping(serve)

	return serve
}
