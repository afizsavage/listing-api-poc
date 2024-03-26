package repository

import (
	"afizsavage/api-poc/entity"
	"time"

	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

type ListingRepository interface {
	Save(listing entity.Listing)
	Update(listing entity.Listing)
	Delete(listing entity.Listing)
	FindAll() []entity.Listing
	GenerateUniqueID() uint64
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

func (db *database) Save(listing entity.Listing) {
	db.connection.Create(&listing)
}

func (db *database) Update(listing entity.Listing) {
	db.connection.Save(&listing)
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
