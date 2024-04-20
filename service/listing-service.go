package service

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/repository"

	"github.com/minio/minio-go/v7"
)

type ListingService interface {
	Save(entity.Listing)  (entity.Listing, error)
	Update(listing entity.Listing) (entity.Listing, error)
	Delete(listing entity.Listing) error
	FindAll() []entity.Listing
	GenerateUniqueID() uint64
	UploadPhoto(uint, *minio.UploadInfo) (entity.Listing, error)
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

func (service *listingService) Save(listing entity.Listing) (entity.Listing, error) {
    savedListing, err := service.listingRepository.Save(listing)
    if err != nil {
        return entity.Listing{}, err
    }
    return savedListing, nil // Return nil for the error here
}

func (service *listingService) UploadPhoto(id uint, uploadedInfo *minio.UploadInfo ) (entity.Listing, error) {
	existingListing, err := service.listingRepository.GetByID(id)

	if err != nil {
		return entity.Listing{}, err
	}

	newPhoto  := entity.Photo {
		Title: uploadedInfo.Key,
		Path: uploadedInfo.Key,
		ListingID: string(rune(existingListing.ID)) ,
	}

	existingListing.Photos = append(existingListing.Photos, newPhoto)
	updatedListing, err := service.listingRepository.Update(existingListing)

	if err != nil {
		return entity.Listing{}, err
	}

	return updatedListing, nil	
}


func (service *listingService) Update(listing entity.Listing) (entity.Listing, error) {
    // Convert listing.ID to uint64
    id := listing.ID

    existingListing, err := service.listingRepository.GetByID(id)
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

    updatedListing, err := service.listingRepository.Update(existingListing)
    if err != nil {
        return entity.Listing{}, err
    }

    return updatedListing, nil
}


func (service *listingService) Delete(listing entity.Listing)   error{
	err := service.listingRepository.Delete(listing)

	return err
}

func (service *listingService) FindAll() []entity.Listing {
	return service.listingRepository.FindAll()
}

