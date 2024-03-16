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

	server.GET("/listings/generate_unique_id", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.GenerateUniqueID())
	})

	server.GET("/listings", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.FindAll())
	})
	
	server.POST("/listings", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.Save(ctx))
	})

	server.PUT("/listings/:id", func(ctx *gin.Context) {
		propertyController.Update(ctx)
	})

	server.DELETE("/listings/:id", func(ctx *gin.Context) {
		propertyController.Delete(ctx)
	})

	server.Run()
}
