package controller

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListingController interface {
	FindAll() []entity.Listing
	Save(ctx *gin.Context) entity.Listing
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type controller struct {
	service service.ListingService
}

func New(service service.ListingService) ListingController {
	return &controller{
		service: service,
	}
}

func (c *controller) FindAll() []entity.Listing {
	return c.service.FindAll()
}

func (c *controller) Save(ctx *gin.Context) entity.Listing {
	var listing entity.Listing
	ctx.BindJSON(&listing)
	c.service.Save(listing)

	return listing
}

func (c *controller) Update(ctx *gin.Context) {
    idStr := ctx.Param("id") // Extracting the ID from the URI
    convertedID, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        // Handle the error (e.g., return an error response)
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var listing entity.Listing
    ctx.BindJSON(&listing)
    // Set the ID of the listing before updating
    listing.ID = convertedID

    c.service.Update(listing)
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
    // No need to bind the JSON to the listing variable if you're not using it

    // Set the ID of the listing before deleting
    listing.ID = convertedID

    c.service.Delete(listing)
}
