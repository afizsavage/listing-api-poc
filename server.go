package main

import (
	"afizsavage/api-poc/controller"
	"afizsavage/api-poc/repository"
	"afizsavage/api-poc/service"

	"github.com/gin-gonic/gin"
)

var (
	listingRepository repository.ListingRepository = repository.NewListingRepository()
	listingService    service.ListingService       = service.New(listingRepository)
	propertyController controller.ListingController = controller.New(listingService)
)

func main() {
	server := gin.Default()

	server.GET("/properties", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.FindAll())
	})

	server.POST("/properties", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.Save(ctx))
	})

	server.PUT("/properties/:id", func(ctx *gin.Context) {
		propertyController.Update(ctx)
	})

	server.DELETE("/properties/:id", func(ctx *gin.Context) {
		propertyController.Delete(ctx)
	})

	server.Run()
}
