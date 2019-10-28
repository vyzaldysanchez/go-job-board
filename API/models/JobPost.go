package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// JobPost represents a job post
type JobPost struct {
	gorm.Model
	UserID      uint      `gorm:"not_null" json:"userId"`
	Title       string    `gorm:"not_null" json:"title"`
	Location    *Location `json:"location,omitempty"`
	LocationID  uint      `gorm:"not_null" json:"locationId"`
	Category    *Category `json:"category,omitempty"`
	CategoryID  uint      `gorm:"not_null" json:"categoryId"`
	Description string    `gorm:"not_null" json:"description"`
	ApplyAt     string    `gorm:"not_null" json:"applyAt"`
	Skills      []Skill   `gorm:"many2many:job_post_skills;" json:"skills,omitempty"`
}

type JobPostService interface {
	JobPostDB
}

type jobPostService struct {
	JobPostDB
}

func NewJobPostService(db *gorm.DB) JobPostService {
	return &jobPostService{
		JobPostDB: &jobPostValidator{
			&jobPostGorm{db},
		},
	}
}

type JobPostDB interface {
	FindAll() ([]JobPost, error)
	FindByParameters(params *JobPost) ([]JobPost, error)
	ByUserID(id uint) ([]JobPost, error)
	ByID(id uint) (*JobPost, error)
	Create(jobPost *JobPost) error
	Update(jobPost *JobPost) error
	Delete(id uint) error
}

type jobPostValidator struct {
	JobPostDB
}

func (jpv *jobPostValidator) Create(jobPost *JobPost) error {

	err := runJobPostValFuncs(
		jobPost, jpv.userIDRequired, jpv.titleRequired, jpv.locationIDRequired, jpv.categoryIDRequired, jpv.descriptionRequired, jpv.applyAtRequired)

	if err != nil {
		return err
	}
	return jpv.JobPostDB.Create(jobPost)
}

func (jpv *jobPostValidator) Update(jobPost *JobPost) error {

	err := runJobPostValFuncs(
		jobPost, jpv.userIDRequired, jpv.titleRequired, jpv.locationIDRequired, jpv.categoryIDRequired, jpv.descriptionRequired, jpv.applyAtRequired)
	if err != nil {
		return err
	}
	return jpv.JobPostDB.Update(jobPost)
}

func (jpv *jobPostValidator) Delete(id uint) error {

	if id <= 0 {
		return ErrIDInvalid
	}

	return jpv.JobPostDB.Delete(id)
}

func (jpv *jobPostValidator) userIDRequired(jp *JobPost) error {
	if jp.UserID <= 0 {
		return ErrUserIDRequired
	}

	return nil
}
func (jpv *jobPostValidator) locationIDRequired(jp *JobPost) error {
	if jp.LocationID <= 0 {
		return ErrLocationIDRequired
	}

	return nil
}

func (jpv *jobPostValidator) categoryIDRequired(jp *JobPost) error {
	if jp.CategoryID <= 0 {
		return ErrCategoryIDRequired
	}

	return nil
}

func (jpv *jobPostValidator) titleRequired(jp *JobPost) error {
	if jp.Title == "" {
		return ErrTitleRequired
	}

	return nil
}

func (jpv *jobPostValidator) descriptionRequired(jp *JobPost) error {
	if jp.Description == "" {
		return ErrDescriptionRequired
	}

	return nil
}

func (jpv *jobPostValidator) applyAtRequired(jp *JobPost) error {
	if jp.ApplyAt == "" {
		return ErrApplyAtRequired
	}

	return nil
}

var _ JobPostDB = &jobPostGorm{}

type jobPostGorm struct {
	db *gorm.DB
}

func (jpg *jobPostGorm) FindAll() ([]JobPost, error) {
	var jobPosts []JobPost
	err := jpg.db.Set("gorm:auto_preload", true).Find(&jobPosts).Error
	if err != nil {
		return nil, err
	}

	return jobPosts, nil
}

// Create will create the provided jobPost and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (jpg *jobPostGorm) Create(jobPost *JobPost) error {
	return jpg.db.Set("gorm:association_autoupdate", false).Create(jobPost).Error
}

func (jpg *jobPostGorm) Update(jobPost *JobPost) error {
	return jpg.db.Save(jobPost).Error
}

func (jpg *jobPostGorm) Delete(id uint) error {
	jobPost := JobPost{Model: gorm.Model{ID: id}}
	return jpg.db.Delete(&jobPost).Error
}

func (jpg *jobPostGorm) ByID(id uint) (*JobPost, error) {
	var jobPost JobPost
	db := jpg.db.Where("id = ?", id)
	err := first(db, &jobPost)

	return &jobPost, err

}

func (jpg *jobPostGorm) FindByParameters(params *JobPost) ([]JobPost, error) {

	var jobPosts []JobPost
	db := jpg.db.Set("gorm:auto_preload", true)
	db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", params.Title))
	params.Title = ""

	err := db.Where(params).Find(&jobPosts).Error

	return jobPosts, err
}

func (jpg *jobPostGorm) ByUserID(id uint) ([]JobPost, error) {
	var jobPosts []JobPost
	err := jpg.db.Where("user_id = ?", id).Find(&jobPosts).Error
	if err != nil {
		return nil, err
	}

	return jobPosts, nil
}

type jobPostValFunc func(*JobPost) error

func runJobPostValFuncs(jobPost *JobPost, fns ...jobPostValFunc) error {
	for _, fn := range fns {
		if err := fn(jobPost); err != nil {
			return err
		}
	}

	return nil
}
