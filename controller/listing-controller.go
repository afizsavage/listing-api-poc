package controller

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/service"

	"github.com/gin-gonic/gin"
)

type ListingController interface {
	FindAll() []entity.Listing
	Save(ctx *gin.Context) entity.Listing
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
