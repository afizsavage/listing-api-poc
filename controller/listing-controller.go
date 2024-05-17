package controller

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/mdf"
	"afizsavage/api-poc/service"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
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
	DeletePhoto(ctx *gin.Context)
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

	err := ctx.Request.ParseMultipartForm(10 << 20) // Parse up to 10 MB of form data
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Listing ID"})
		return
	}

	id := uint(id64)

	form := ctx.Request.MultipartForm
	if form == nil || form.File == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No photo provided"})
		return
	}

	files := form.File["photos"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No photo provided"})
		return
	}

	var uploadedFiles []entity.Photo


	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
			return
		}
		defer file.Close()

		fileName := fileHeader.Filename
		fileSize := fileHeader.Size

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
			return
		}

		detectedFileType := http.DetectContentType(fileBytes)
		parsedFT := strings.Split(detectedFileType, ";")

		objectName := "listing-photos/" + fileName

		// Upload the file to your storage (e.g., MinIO)
		uploadedPhoto, err := minioClient.PutObject(
			ctx,
			bucketName,
			objectName,
			bytes.NewReader(fileBytes),
			fileSize,
			minio.PutObjectOptions{ContentType: parsedFT[0], PartSize: uint64(partSize)})
		if err != nil {
			fmt.Println("Upload photo error:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't upload file"})
			return
		}

		newPhoto  := entity.Photo {
			Title: filepath.Base(uploadedPhoto.Key),
			Path: uploadedPhoto.Key,
			ListingID: id ,
		}

		uploadedFiles = append(uploadedFiles, newPhoto)
	}

	// Process the uploaded files (e.g., update a listing with the photos)
    updatedListing, err := c.service.UploadPhotos(id, uploadedFiles)
    if err != nil {
        fmt.Println("Update listing error:", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload photos"})
        return
    }

    // Return a success message indicating all photos were uploaded successfully
    ctx.JSON(http.StatusOK, gin.H{"message": "Photos updated successfully", "data": updatedListing})

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

func (c *controller) DeletePhoto(ctx *gin.Context) {
	idStr := ctx.Param("id")
	photoIDStr := ctx.Param("photo-id")


    listingID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Listing ID"})
        return
    }

	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Photo ID"})
        return
    }

	var photo entity.Photo

	photo.ID = uint(photoID)
	
	updatedListing, err := c.service.DeletePhoto(uint(listingID),photo)
	
	
	if err != nil {
		fmt.Println("update listing error:", err) // Print the error message for debugging

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed delete photo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Photo  successfully", "data": updatedListing})

}
