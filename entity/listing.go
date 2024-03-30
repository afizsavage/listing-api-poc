package entity

import (
	"github.com/shopspring/decimal"
)

type Listing struct {
	ID         uint64  `json:"id" gorm:"primary_key;auto_increment"`
	Country    string  `json:"country" gorm:"type:varchar(200)"`
	City       string  `json:"city" gorm:"type:varchar(200)"`
	Address    string  `json:"address" gorm:"type:varchar(200)"`
	Bedrooms   string  `json:"bedrooms" gorm:"type:varchar(50)"`
	Bathrooms  string  `json:"bathrooms" gorm:"type:varchar(50)"`
	Type       string  `json:"type" gorm:"type:varchar(50)"`
	Title      string  `json:"title" gorm:"type:varchar(50)"`
	Latitude   decimal.Decimal `json:"latitude" gorm:"type:decimal(20,8);"`
	Longitude  decimal.Decimal `json:"longitude" gorm:"type:decimal(20,8);"`
	ExternalID  uint64  `json:"external_d" gorm:"type:bigint"`
}