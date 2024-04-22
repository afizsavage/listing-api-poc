package main

import (
	"afizsavage/api-poc/controller"
	"afizsavage/api-poc/mdf"
	"afizsavage/api-poc/repository"
	"afizsavage/api-poc/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


var (
	listingRepository repository.ListingRepository = repository.NewListingRepository()
	listingService    service.ListingService       = service.New(listingRepository)
	propertyController controller.ListingController = controller.New(listingService)
)

func main() {
	server := gin.Default()

	server.Use(cors.Default())
	minioClient := mdf.InitMinioClient()

	log.Printf("%#v\n", minioClient)
	
	server.GET("/listings/generate_unique_id", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.GenerateUniqueID())
	})

	server.GET("/listings", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.FindAll())
	})

	server.GET("/listings/:id", func(ctx *gin.Context) {
		 propertyController.GetByID(ctx)
	})

	server.GET("/listings/photos/:photo-name", func(ctx *gin.Context) {
		propertyController.GetImageURL(ctx)
	})

	server.POST("/listings/:id/upload/photo", func(ctx *gin.Context) {
		propertyController.UploadPhoto(ctx)
	})
	
	server.POST("/listings", func(ctx *gin.Context) {
		 propertyController.Save(ctx)
	})

	server.PUT("/listings/:id", func(ctx *gin.Context) {
		propertyController.Update(ctx)
	})

	server.DELETE("/listings/:id", func(ctx *gin.Context) {
		propertyController.Delete(ctx)
	})

	server.Run()
}
