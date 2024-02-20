package main

import (
	"afizsavage/api-poc/controller"
	"afizsavage/api-poc/service"

	"github.com/gin-gonic/gin"
)

var (
	propertyService service.ListingService = service.New()
	propertyController controller.ListingController = controller.New(propertyService)
)

func main() {
	server := gin.Default()

	server.GET("/properties", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.FindAll())
	})

	server.POST("/properties", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.Save(ctx))
	})

	server.Run()
}