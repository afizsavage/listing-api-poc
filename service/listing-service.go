package service

import (
	"afizsavage/api-poc/entity"
	"afizsavage/api-poc/repository"
	"fmt"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

type ListingService interface {
	Save(entity.Listing)  (entity.Listing, error)
	Update(listing entity.Listing) (entity.Listing, error)
	Delete(listing entity.Listing) error
	DeletePhoto(listingID uint, photo entity.Photo ) (entity.Listing, error)
	GetAll() []entity.Listing
	GenerateUniqueID() uint64
	UploadPhoto(uint, *minio.UploadInfo) (entity.Listing, error)
	GetByID(uint) (entity.Listing, error)
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
		Title: filepath.Base(uploadedInfo.Key),
		Path: uploadedInfo.Key,
		ListingID: existingListing.ID ,
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




func (service *listingService) GetAll() []entity.Listing {
	listings, err :=  service.listingRepository.GetAll()
	
	if err != nil {
        fmt.Println("find all listing error", err)
    }

	return listings
}

func (service *listingService) GetByID(id uint) (entity.Listing, error) {
	listing, err := service.listingRepository.GetByID(id)

	if err != nil {
        fmt.Println("get listing by id err", err)
    }

	return listing, err

}

func (service *listingService) Delete(listing entity.Listing)   error{
	err := service.listingRepository.Delete(listing)

	return err
}

func (service *listingService) DeletePhoto(listingID uint, photo entity.Photo)  ( entity.Listing,error){

	err := service.listingRepository.DeletePhoto(photo)
	if err != nil {
        return entity.Listing{}, err
    }

	updatedListing, err := service.listingRepository.GetByID(listingID)
	if err != nil {
        return entity.Listing{}, err
    }

	return updatedListing, err	
}