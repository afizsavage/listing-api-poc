package repository

import (
	"afizsavage/api-poc/entity"

	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ListingRepository interface {
	Save(listing entity.Listing) entity.Listing
	Update(listing entity.Listing) (entity.Listing, error)
	Delete(listing entity.Listing)
	FindAll() []entity.Listing
	GenerateUniqueID() uint64
	GetByID(id uint64) (entity.Listing, error)
}

type database struct {
	connection *gorm.DB
}

func NewListingRepository() ListingRepository {
	dsn := "host=127.0.0.1 user=postgres password=password dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&entity.Listing{})

	return &database{
		connection: db,
	}
}

func (db *database) Save(listing entity.Listing) entity.Listing {
    db.connection.Create(&listing)
    return listing
}

// Update updates a listing in the database.
func (db *database) Update(listing entity.Listing) (entity.Listing, error) {
	// Perform database update operation
	if err := db.connection.Save(&listing).Error; err != nil {
		return entity.Listing{}, err
	}
	return listing, nil
}


func (db *database) Delete(listing entity.Listing) {
	db.connection.Delete(&listing)
}

func (db *database) FindAll() []entity.Listing {
	var listings []entity.Listing
	db.connection.Find(&listings)
	return listings
}

func (db *database) GenerateUniqueID() uint64 {
    // Get the current timestamp in Unix format and return it as uint64
    return uint64(time.Now().Unix())
}

func (db *database) GetByID(id uint64) (entity.Listing, error) {
    var listing entity.Listing
    if err := db.connection.First(&listing, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return entity.Listing{}, errors.New("listing not found")
        }
        return entity.Listing{}, err
    }
    return listing, nil
}
