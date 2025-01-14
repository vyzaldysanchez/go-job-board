package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		//db.LogMode(true)
		s.db = db
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {

	return func(s *Services) error {

		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithJobPost() ServicesConfig {

	return func(s *Services) error {
		s.JobPost = NewJobPostService(s.db)
		return nil
	}
}

func WithCategory() ServicesConfig {
	return func(s *Services) error {
		s.Category = NewCategoryService(s.db)
		return nil
	}
}
func WithLocation() ServicesConfig {
	return func(s *Services) error {
		s.Location = NewLocationService(s.db)
		return nil
	}
}

func WithSkill() ServicesConfig {

	return func(s *Services) error {
		s.Skill = NewSkillService(s.db)
		return nil
	}
}

func WithOAuth() ServicesConfig {
	return func(s *Services) error {

		s.OAuth = NewOAuthService(s.db)
		return nil
	}
}
func WithLogMode(logMode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(logMode)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}

	return &s, nil

}

type Services struct {
	JobPost  JobPostService
	Category CategoryService
	Location LocationService
	User     UserService
	Skill    SkillsService
	OAuth    OAuthService
	db       *gorm.DB
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will attempt to automatically migrate
// all tables
func (s *Services) AutoMigrate() error {
	err := s.db.AutoMigrate(
		&User{},
		&Role{},
		&Location{},
		&Category{},
		&JobPost{},
		&Skill{},
		&CompanyProfile{},
		&CompanyBenefit{},
		&pwReset{},
		&OAuth{}).Error
	if err != nil {
		return err
	}
	return runPopulatingFuncs(s.seedRoles, s.seedLocations, s.seedCategories, s.seedSkills)
}
func (s *Services) seedRoles() error {
	return s.db.Model(&Role{}).Create(&Role{RoleName: "User"}).Create(&Role{RoleName: "Candidate"}).Error
}

func (s *Services) seedLocations() error {
	db := s.db.Model(&Location{})
	locationsSeed := s.GetLocationsSeed()
	for _, c := range locationsSeed {
		db = db.Create(&c)
	}
	return db.Error
}

func (s *Services) GetLocationsSeed() []Location {
	return []Location{
		Location{LocationName: "USA"},
		Location{LocationName: "Canada"},
		Location{LocationName: "Europe"},
		Location{LocationName: "Remote"},
	}
}

func (s *Services) GetCategoriesSeed() []Category {
	return []Category{
		Category{CategoryName: "Web Development"},
		Category{CategoryName: "Mobile Development"},
		Category{CategoryName: "QA"},
		Category{CategoryName: "DBA"},
		Category{CategoryName: "DevOps"},
	}
}
func (s *Services) seedCategories() error {
	db := s.db.Model(&Category{})
	categoriesSeed := s.GetCategoriesSeed()
	for _, c := range categoriesSeed {
		db = db.Create(&c)
	}
	return db.Error
}
func (s *Services) seedSkills() error {
	return s.db.Model(&Skill{}).
		Create(&Skill{SkillName: "JavaScript"}).
		Create(&Skill{SkillName: "Golang"}).Error
}

// DestructiveReset drops the all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.Exec("DROP TABLE IF EXISTS job_post_skills;").DropTableIfExists(
		&User{},
		&Role{},
		&JobPost{},
		&Category{},
		&Location{},
		&Skill{},
		&CompanyProfile{},
		&CompanyBenefit{},
		&pwReset{},
		&OAuth{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

type populatingFunc func() error

func runPopulatingFuncs(fns ...populatingFunc) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
