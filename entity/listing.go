package entity

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Photo struct {
	gorm.Model
	Title     string `json:"Title" gorm:"type:varchar(250)"`
	Path      string `json:"Path" gorm:"type:varchar(250)"`
	ListingID uint   
}
type Listing struct {
	gorm.Model
	Country    string         `json:"Country" gorm:"type:varchar(200)"`
	City       string         `json:"City" gorm:"type:varchar(200)"`
	Address    string         `json:"Address" gorm:"type:varchar(200)"`
	Bedrooms   string         `json:"Bedrooms" gorm:"type:varchar(50)"`
	Bathrooms  string         `json:"Bathrooms" gorm:"type:varchar(50)"`
	Type       string         `json:"Type" gorm:"type:varchar(50)"`
	Title      string         `json:"Title" gorm:"type:varchar(50)"`
	Latitude   decimal.Decimal`json:"Latitude" gorm:"type:decimal(20,8);"`
	Longitude  decimal.Decimal`json:"Longitude" gorm:"type:decimal(20,8);"`
	ExternalID uint64         `json:"External_d" gorm:"type:bigint"`
	Amenities  pq.StringArray `gorm:"type:text[]"`
	Photos     []Photo
}
