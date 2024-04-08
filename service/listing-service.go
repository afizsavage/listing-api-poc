package service

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/repository"
	"errors"
)

type ListingService interface {
	Save(entity.Listing) entity.Listing
	Update(listing entity.Listing) (entity.Listing, error)
	Delete(listing entity.Listing)
	FindAll() []entity.Listing
	GenerateUniqueID() uint64
}

type listingService struct {
	listingRepository repository.ListingRepository
}

func New(repo repository.ListingRepository) ListingService {
	return &listingService{
		listingRepository: repo,
	}
}

func (service *listingService) GenerateUniqueID() uint64 {
	return service.listingRepository.GenerateUniqueID()
}

func (service *listingService) Save(listing entity.Listing) entity.Listing {
    savedListing := service.listingRepository.Save(listing)
    return savedListing
}

func (s *listingService) Update(listing entity.Listing) (entity.Listing, error) {
	if listing.ID == 0 {
		return entity.Listing{}, errors.New("ID is required for updating a listing")
	}

	existingListing, err := s.listingRepository.GetByID(listing.ID) // Get the existing listing from the repository
	if err != nil {
		return entity.Listing{}, err
	}

	// Update the fields of the existing listing with new values
	existingListing.Country = listing.Country
	existingListing.City = listing.City
	existingListing.Address = listing.Address
	existingListing.Bedrooms = listing.Bedrooms
	existingListing.Bathrooms = listing.Bathrooms
	existingListing.Type = listing.Type
	existingListing.Title = listing.Title
	existingListing.Latitude = listing.Latitude
	existingListing.Longitude = listing.Longitude
	existingListing.ExternalID = listing.ExternalID

	updatedListing, err := s.listingRepository.Update(existingListing)
	if err != nil {
		return entity.Listing{}, err
	}
	return updatedListing, nil
}


func (service *listingService) Delete(listing entity.Listing) {
	service.listingRepository.Delete(listing)
}

func (service *listingService) FindAll() []entity.Listing {
	return service.listingRepository.FindAll()
}

