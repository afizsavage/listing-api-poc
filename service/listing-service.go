package service

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/repository"
)

type ListingService interface {
	Save(entity.Listing) entity.Listing
	Update(listing entity.Listing)
	Delete(listing entity.Listing)
	FindAll() []entity.Listing
	GenerateUniqueID() string
}

type listingService struct {
	listingRepository repository.ListingRepository
}

func New(repo repository.ListingRepository) ListingService {
	return &listingService{
		listingRepository: repo,
	}
}

func (service *listingService) GenerateUniqueID() string {
	return service.listingRepository.GenerateUniqueID()
}

func (service *listingService) Save(listing entity.Listing) entity.Listing {
    savedListing := service.listingRepository.Save(listing)
    return savedListing
}

func (service *listingService) Update(listing entity.Listing) {
	service.listingRepository.Update(listing)
}

func (service *listingService) Delete(listing entity.Listing) {
	service.listingRepository.Delete(listing)
}

func (service *listingService) FindAll() []entity.Listing {
	return service.listingRepository.FindAll()
}
