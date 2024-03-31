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
	Save(ctx *gin.Context) 
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
    GenerateUniqueID() string
}

type controller struct {
	service service.ListingService
}

func New(service service.ListingService) ListingController {
	return &controller{
		service: service,
	}
}

func (c *controller) GenerateUniqueID() string {
	return c.service.GenerateUniqueID()
}

func (c *controller) FindAll() []entity.Listing {
	return c.service.FindAll()
}

func (c *controller) Save(ctx *gin.Context) {
    var listing entity.Listing
    ctx.BindJSON(&listing)
    createdListing := c.service.Save(listing)

    ctx.JSON(http.StatusCreated, gin.H{"message": "Listing created successfully", "listing": createdListing})
}

// Update updates a listing based on the ID provided in the URL path.
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

	listing.ID = id // Set the ID from the URL path

	updatedListing, err := c.service.Update(listing)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update listing"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Listing updated successfully", "listing": updatedListing})
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
