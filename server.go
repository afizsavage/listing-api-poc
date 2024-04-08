package main

import (
	"context"
	"log"

	"afizsavage/api-poc/controller"
	"afizsavage/api-poc/repository"
	"afizsavage/api-poc/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	listingRepository repository.ListingRepository = repository.NewListingRepository()
	listingService    service.ListingService       = service.New(listingRepository)
	propertyController controller.ListingController = controller.New(listingService)
)

func main() {
	server := gin.Default()

	server.Use(cors.Default())

	endpoint := "127.0.0.1:80"
	accessKeyID := "lOGD8QU43Y3ElIB4"
	secretAccessKey := "YA9bL6RZTCDBSGp7Zh6RtmXFu477wjfb"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	buckets, err :=  minioClient.ListBuckets(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("TOTAL BUCKETS", len(buckets)) // minioClient is now setup

	server.GET("/listings/generate_unique_id", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.GenerateUniqueID())
	})

	server.GET("/listings", func(ctx *gin.Context) {
		ctx.JSON(200, propertyController.FindAll())
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
