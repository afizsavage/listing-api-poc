package entity

type Listing struct {
	ID uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Country string `json:"country" gorm:"type:varchar(200)"`
	City string `json:"city" gorm:"type:varchar(200)"`
	Address string `json:"address" gorm:"type:varchar(200)"`
}