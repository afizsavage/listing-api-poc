package service

import (
	"afizsavage/api-poc/entity"
)

type ListingService interface {
	Save(entity.Listing) entity.Listing
	FindAll() []entity.Listing
}

type listingService struct {
	listings []entity.Listing
}

func New() ListingService {
	return &listingService{}
}

func (service *listingService) Save(property entity.Listing) entity.Listing {
	service.listings = append(service.listings, property)

	return property
}

func (service *listingService)	FindAll() []entity.Listing {
	return service.listings
}