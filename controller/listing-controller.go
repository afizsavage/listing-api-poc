package controller

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/mdf"
	"afizsavage/api-poc/service"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type ListingController interface {
	GetAll() []entity.Listing
	Save(ctx *gin.Context) 
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
    GenerateUniqueID() uint64
	UploadPhoto(ctx *gin.Context)
	GetImageURL(ctx *gin.Context)
	GetByID(ctx *gin.Context) entity.Listing

}

type controller struct {
	service service.ListingService
}

const bucketName = "propati-poc"
const partSize int = 20 * 1024 * 1024
const directory = "listing-photos"
 
var minioClient = mdf.InitMinioClient()


func New(service service.ListingService) ListingController {
	return &controller{
		service: service,
	}
}

func (c *controller) GenerateUniqueID() uint64 {
	return c.service.GenerateUniqueID()
}

func (c *controller) GetAll() []entity.Listing {
	return c.service.GetAll()
}

func (c *controller) GetImageURL(ctx *gin.Context)  {
	photoName := ctx.Param("photo-name")
	objectName := directory + "/" + photoName

	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
    	fmt.Println(err)
    	return
	}

	defer object.Close()

	ctx.Header("Content-Type", "image/jpeg")

    // Stream the image data from Minio to the HTTP response
    if _, err := io.Copy(ctx.Writer, object); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream image"})
        return
    }
}


func (c *controller) GetByID(ctx *gin.Context) entity.Listing  {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return entity.Listing{}
	}

	id := uint(id64)

	returnedListing, err := c.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return entity.Listing{}
	}


	ctx.JSON(http.StatusOK, gin.H{ "data": returnedListing})


	return returnedListing
	
}

func (c *controller) Save(ctx *gin.Context) {
    var listing entity.Listing
    ctx.BindJSON(&listing)
    createdListing, err := c.service.Save(listing)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cant create listing"})
		return
	}

    ctx.JSON(http.StatusCreated, gin.H{"message": "Listing created successfully", "data": createdListing})
}

func (c *controller) UploadPhoto(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	id := uint(id64)

	file, fileHeader, err := ctx.Request.FormFile("photo")

	if err != nil {
		   ctx.JSON(http.StatusBadRequest, gin.H{"error": "No photo provided"})
		   return
	}

	defer file.Close()
   
	fileName := fileHeader.Filename
	fileSize := fileHeader.Size

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	detectedFileType := http.DetectContentType(fileBytes)
	parsedFT := strings.Split(detectedFileType, ";")

	objectName := "listing-photos/" + fileName

 
	uploadedPhoto, err := minioClient.PutObject(
			context.Background(),  
			bucketName,
		 	objectName,
		  	bytes.NewReader(fileBytes), 
		  	fileSize,
			minio.PutObjectOptions{ContentType: parsedFT[0], PartSize: uint64(partSize)})

	if err != nil {
		fmt.Println("upload photo error:", err) 
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't upload file"})
	}

	updatedListing, err := c.service.UploadPhoto(id, &uploadedPhoto )
	
	
	if err != nil {
		fmt.Println("update listing error:", err) // Print the error message for debugging

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload photo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Photo updated successfully", "data": updatedListing})
}


func (c *controller) Update(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var listing entity.Listing
    if err := ctx.ShouldBindJSON(&listing); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Convert id to uint and assign it to listing.ID
    listing.ID = uint(id)

    updatedListing, err := c.service.Update(listing)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update listing"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Listing updated successfully", "data": updatedListing})
}

func (c *controller) Delete(ctx *gin.Context) {
    idStr := ctx.Param("id") // Extracting the ID from the URI
    convertedID, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        // Handle the error (e.g., return an error response)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var listing entity.Listing

    // Convert convertedID to uint and assign it to listing.ID
    listing.ID = uint(convertedID)

    c.service.Delete(listing)
}
