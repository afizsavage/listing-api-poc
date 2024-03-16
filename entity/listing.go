package entity

type Listing struct {
	ID uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Country string `json:"country" gorm:"type:varchar(200)"`
	City string `json:"city" gorm:"type:varchar(200)"`
	Address string `json:"address" gorm:"type:varchar(200)"`
	Bedrooms string `json:"bedrooms" gorm:"type:varchar(50)"`
	Bathrooms string  `json:"bathrooms" gorm:"type:varchar(50)"`
	Type string  `json:"type" gorm:"type:varchar(50)"`
	Title string  `json:"title" gorm:"type:varchar(50)"`
	Latitude string  `json:"latitude" gorm:"type:varchar(50)"`
	Longitude string  `json:"longitude" gorm:"type:varchar(50)"`
}