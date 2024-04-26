package repository

import (
	"afizsavage/api-poc/entity"
	"fmt"

	"errors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ListingRepository interface {
	Save(listing entity.Listing) (entity.Listing, error )
	Update(listing entity.Listing) (entity.Listing, error)
	Delete(listing entity.Listing) error
	GetAll() ([]entity.Listing, error)
	GenerateUniqueID() uint64
	GetByID(id uint) (entity.Listing, error)
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

	db.AutoMigrate(&entity.Photo{})
	db.AutoMigrate(&entity.Listing{})

	return &database{
		connection: db,
	}
}

func (db *database) Save(listing entity.Listing) (entity.Listing, error) {
    if err := db.connection.Create(&listing).Error; err != nil {
        fmt.Println("Database error:", err) // Print the error message for debugging
        return entity.Listing{}, err
    }
    return listing, nil
}

func (db *database) Delete(listing entity.Listing) error {
    if err := db.connection.Delete(&listing).Error; err != nil {
        return err
    }
    return nil
}
// Update updates a listing in the database.
func (db *database) Update(listing entity.Listing) (entity.Listing, error) {
	// Perform database update operation
	if err := db.connection.Save(&listing).Error; err != nil {
		fmt.Println("update database error", err)

		return entity.Listing{}, err
	}
	return listing, nil
}

func (db *database) GetAll() ([]entity.Listing, error) {
	var listings []entity.Listing
    err := db.connection.Model(&entity.Listing{}).Preload("Photos").Find(&listings).Error
    
	return listings, err
}

func (db *database) GetByID(id uint) (entity.Listing, error) {
    var listing entity.Listing
	
    if err := db.connection.Model(&entity.Listing{}).Preload("Photos").First(&listing, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return entity.Listing{}, errors.New("listing not found")
        }
        return entity.Listing{}, err
    }

    return listing, nil
}

func (db *database) GenerateUniqueID() uint64 {
    // Get the current timestamp in Unix format and return it as uint64
    return uint64(time.Now().Unix())
}
