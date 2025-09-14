package main

import (
	ginadapter "github.com/DrWeltschmerz/users-adapter-gin/ginadapter"
	core "github.com/DrWeltschmerz/users-core"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Users API
// @version 1.0
// @description This is a sample server for user management.
// @BasePath /
//
// NOTE: This main.go is for demo and Swagger preview only.
// Real wiring of repositories, hasher, and tokenizer should be done in your application, as documented in users-core.
func main() {
	r := gin.Default()

	// Swagger docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Demo only: do not use nils in production!
	var svc *core.Service = nil
	var tokenizer core.Tokenizer = nil
	ginadapter.RegisterRoutes(r, svc, tokenizer)

	r.Run()
}
