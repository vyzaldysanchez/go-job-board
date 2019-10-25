package models

import "github.com/jinzhu/gorm"

type Location struct {
	gorm.Model
	LocationName string `json:"locationName"`
}

type LocationService interface {
	LocationDB
}

type locationService struct {
	LocationDB
}

func NewLocationService(db *gorm.DB) LocationService {
	return &locationService{
		LocationDB: &locationGorm{db},
	}
}

type locationGorm struct {
	db *gorm.DB
}

func (lg *locationGorm) FindAll() ([]Location, error) {
	var locations []Location
	err := lg.db.Find(&locations).Error
	if err != nil {
		return nil, err
	}

	return locations, nil
}

type LocationDB interface {
	FindAll() ([]Location, error)
}
