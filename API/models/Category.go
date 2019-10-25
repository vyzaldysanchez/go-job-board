package models

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model
	CategoryName string `json:"categoryName"`
}

type CategoryService interface {
	CategoryDB
}

type categoryService struct {
	CategoryDB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{
		CategoryDB: &categoryGorm{db},
	}
}

type categoryGorm struct {
	db *gorm.DB
}

func (cg *categoryGorm) FindAll() ([]Category, error) {
	var categories []Category
	err := cg.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

type CategoryDB interface {
	FindAll() ([]Category, error)
}
