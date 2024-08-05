package api

import (
	"user/api/handler"
	_ "user/api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title User Service API
// @version 1.0
// @description This is a sample server for a user service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:50052
// @BasePath /
func NewRouter(hand *handler.Handler) *gin.Engine{
	router:=gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/register",hand.Register)
	router.POST("/login",hand.LoginUser)
	router.POST("/refresh-token",hand.RefReshToken)
	return router
}