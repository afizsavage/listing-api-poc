package entity

type Listing struct {
	ID         uint64 `json:"id" gorm:"primary_key;auto_increment"`
	ExternalID uint64 `json:"external_id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Address    string `json:"address"`
	Bedrooms   uint   `json:"bedrooms"`
	Bathrooms  uint   `json:"bathrooms"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}