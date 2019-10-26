package model_services_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

type PostgressConfigTest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func Dialect() string {
	return "postgres"
}
func DefaultPostgressConfig() PostgressConfigTest {

	return PostgressConfigTest{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Dbname:   "job_board_test",
	}
}

func ConnectionInfo() string {
	c := DefaultPostgressConfig()
	env := os.Getenv("env")
	if env == "CI" {
		c.Password = ""
	}

	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Dbname)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Dbname)
}

func TestUserService(t *testing.T) {

	services, err := models.NewServices(
		models.WithGorm(
			Dialect(),
			ConnectionInfo()),
		models.WithLogMode(false),
		models.WithUser("pepperhere", "randomtesthmacvalue"),
		models.WithSkill(),
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())
	t.Run("Create", testUserService_Create(services.User))
	t.Run("Find", testUserService_Find(services.User))
	t.Run("Update", testUserService_Update(services.User, services.Skill))
	t.Run("Delete", testUserService_Delete(services.User))

}

func testUserService_Update(us models.UserService, ss models.SkillsService) func(t *testing.T) {
	return func(t *testing.T) {

		t.Run("CompanyProfile", func(t *testing.T) {
			user := findUserByID(us, 1, t)
			companyProfile := testUpdateCompanyProfileFields(user, us, t)
			testAddCompanyProfileSkill(t, ss, companyProfile)
			testRemoveCompanyProfileSkill(t, ss, companyProfile)
			testAddCompanyProfileBenefit(t, us, companyProfile)
			testUpdateCompanyProfileBenefit(t, companyProfile, us)
			testRemoveCompanyProfileBenefit(t, companyProfile, us)
		})

	}
}

func testRemoveCompanyProfileBenefit(t *testing.T, got *models.CompanyProfile, us models.UserService) {
	t.Run("remove-benefit", func(t *testing.T) {
		benefit := got.CompanyBenefits[0]

		if err := us.RemoveCompanyProfileBenefit(got, benefit); err != nil {
			t.Errorf("error removing company profile benefit, error = %v, companyProfile = %v, benefit = %v", err, got, benefit)
		}

		got = findUserByID(us, 1, t).CompanyProfile

		if len(got.CompanyBenefits) != 0 {
			t.Errorf("expected benefits list to be empty, but got = %v elements", len(got.CompanyBenefits))
		}
	})
}

func testUpdateCompanyProfileFields(user *models.User, us models.UserService, t *testing.T) *models.CompanyProfile {
	user.CompanyProfile = &models.CompanyProfile{
		Website:        "samysoft.com",
		FoundedYear:    1991,
		Description:    "a very nice company 2",
		CompanyLogoUrl: "a logo url",
	}
	if err := us.Update(user); err != nil {
		t.Fatal(err)
	}
	got := user.CompanyProfile
	want := &models.CompanyProfile{
		Website:        "samysoft.com",
		FoundedYear:    1991,
		Description:    "a very nice company 2",
		CompanyLogoUrl: "a logo url",
	}
	if got.Website != want.Website {
		t.Errorf("Website did not update correctly got = %v, want = %v", got.Website, want.Website)
	}
	if got.FoundedYear != want.FoundedYear {
		t.Errorf("FoundedYear did not update correctly got = %v, want = %v", got.FoundedYear, want.FoundedYear)
	}
	if got.Description != want.Description {
		t.Errorf("Description did not update correctly got = %v, want = %v", got.Description, want.Description)
	}
	if got.CompanyLogoUrl != want.CompanyLogoUrl {
		t.Errorf("Description did not update correctly got = %v, want = %v", got.CompanyLogoUrl, want.CompanyLogoUrl)
	}
	return got
}

func testUpdateCompanyProfileBenefit(t *testing.T, got *models.CompanyProfile, us models.UserService) {
	t.Run("update-benefit", func(t *testing.T) {
		benefit := &got.CompanyBenefits[0]
		benefit.BenefitName = "Remote Work"

		if err := us.UpdateCompanyProfileBenefit(benefit); err != nil {
			t.Errorf("error updating company profile benefit, error = %v, companyProfile = %v, benefit = %v", err, got, benefit)
		}

		got = findUserByID(us, 1, t).CompanyProfile

		if got.CompanyBenefits[0].BenefitName != benefit.BenefitName {
			t.Errorf("expected added benefit name to be %q, but got %q", benefit.BenefitName, got.CompanyBenefits[0].BenefitName)
		}

		testCompanyProfileIsNotNil(t, us, got.CompanyBenefits[0])
		testCompanyProfileBenefitNameIsNotEmpty(t, us, got)

	})
}

func testAddCompanyProfileBenefit(t *testing.T, us models.UserService, got *models.CompanyProfile) {
	t.Run("add-benefit", func(t *testing.T) {
		benefit := models.CompanyBenefit{
			BenefitName: "Flexible schedule",
		}

		if err := us.AddCompanyProfileBenefit(got, benefit); err != nil {
			t.Errorf("error adding company profile benefit, error = %v, companyProfile = %v, benefit = %v", err, got, benefit)
		}

		got := findUserByID(us, 1, t).CompanyProfile

		if got.CompanyBenefits[0].BenefitName != benefit.BenefitName {
			t.Errorf("expected added benefit name to be %q, but got = %q", benefit.BenefitName, got.CompanyBenefits[0].BenefitName)
		}

		testCompanyProfileIsNotNil(t, us, benefit)

		testCompanyProfileBenefitNameIsNotEmpty(t, us, got)

	})
}

func testCompanyProfileBenefitNameIsNotEmpty(t *testing.T, us models.UserService, companyProfile *models.CompanyProfile) {
	t.Run("SadPath: empty BenefitName is not allowed", func(t *testing.T) {
		wantError := models.ErrBenefitNameRequired
		benefit := models.CompanyBenefit{
		}
		if err := us.AddCompanyProfileBenefit(companyProfile, benefit); err != wantError {
			t.Errorf("should return %q error got %q error", wantError, err)
		}
	})
}

func testCompanyProfileIsNotNil(t *testing.T, us models.UserService, benefit models.CompanyBenefit) {
	t.Run("SadPath: nil company profile not allowed", func(t *testing.T) {
		wantError := models.ErrCompanyProfileRequired
		if err := us.AddCompanyProfileBenefit(nil, benefit); err != wantError {

			t.Errorf("should return %q error got %q error", wantError, err)
		}
	})
}

func testAddCompanyProfileSkill(t *testing.T, ss models.SkillsService, got *models.CompanyProfile) {
	testAddSkill(t, ss, got)
}

func testAddSkill(t *testing.T, ss models.SkillsService, got interface{}) {
	t.Run("add-skill", func(t *testing.T) {
		skill := models.Skill{}
		skill.ID = 1

		if err := ss.AddSkillToOwner(got, skill); err != nil {
			t.Errorf("error adding skill, error = %v, model = %v, skill = %v", err, got, skill)
		}

		t.Run("SadPath: skill with no ID is not allowed", func(t *testing.T) {
			wantError := models.ErrIDInvalid
			if err := ss.AddSkillToOwner(got, models.Skill{}); err != wantError {
				t.Errorf("should return %q error got %q error", wantError, err)
			}
		})

	})
}
func testRemoveCompanyProfileSkill(t *testing.T, ss models.SkillsService, got *models.CompanyProfile) {
	testRemoveSkill(t, ss, got)
	if len(got.Skills) != 0 {
		t.Errorf("expected skills list to be empty, but got = %v elements", len(got.Skills))
	}
}

func testRemoveSkill(t *testing.T, ss models.SkillsService, got interface{}) {
	t.Run("remove-skill", func(t *testing.T) {
		skill := models.Skill{}
		skill.ID = 1

		if err := ss.DeleteSkillFromOwner(got, skill); err != nil {
			t.Errorf("error removing skill, error = %v, model = %v, skill = %v", err, got, skill)
		}

	})
}

func testUserService_Delete(us models.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		if err := us.Delete(1); err != nil {
			t.Error(err)
		}
	}
}

func testUserService_Find(us models.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		want := models.User{
			Email: "ps3_3@hotmail.com",
		}
		t.Run("ByID", func(t *testing.T) {
			want.ID = 1
			got := findUserByID(us, want.ID, t)
			if want.ID != got.ID {
				t.Fatalf("invalid user retrieved, got user with ID:%v, want user with ID %v", got.ID, want.ID)
			}
		})
		t.Run("ByEmail", func(t *testing.T) {
			got, err := us.ByEmail(want.Email)
			if err != nil {
				t.Error(err)
			}
			if want.Email != got.Email {
				t.Fatalf("invalid user retrieved, got user with email:%v, want user with ID %v", got.ID, want.ID)
			}
		})
	}
}

func findUserByID(us models.UserService, id uint, t *testing.T) *models.User {
	got, err := us.ByID(id)
	if err != nil {
		t.Error(err)
	}
	return got
}

func testUserService_Create(us models.UserService) func(t *testing.T) {

	return func(t *testing.T) {
		newUser := models.User{
			Email:    "ps3_3@hotmail.com",
			Password: "megaman007",
		}
		if err := us.Create(&newUser); err != nil {
			t.Error(err)
		}
	}
}

func TestJobsService(t *testing.T) {

	services, err := models.NewServices(
		models.WithGorm(
			Dialect(),
			ConnectionInfo()),
		models.WithLogMode(false),
		models.WithJobPost(),
		models.WithSkill(),
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())
	t.Run("Create", testJobsService_Create(services.JobPost))
	t.Run("Find", testJobsService_Find(services.JobPost))
	t.Run("Update", testJobsService_Update(services.JobPost, services.Skill))
	t.Run("Delete", testJobsService_Delete(services.JobPost))

}

func TestLocationsService(t *testing.T) {

	services, err := models.NewServices(
		models.WithGorm(
			Dialect(),
			ConnectionInfo()),
		models.WithLogMode(false),
		models.WithLocation(),
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())

	t.Run("Find", testLocationsService_Find(services.Location, services.GetLocationsSeed))
}

func testLocationsService_Find(ls models.LocationService, getLocationsSeed func() []models.Location) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("FindAll", func(t *testing.T) {
			want := len(getLocationsSeed())
			got, err := ls.FindAll()
			if err != nil {
				t.Fatalf("locations could not be fetched error: %s", err.Error())
			}
			if len(got) <= 0 {
				t.Fatalf("expected to find %d locations but got %d locations", want, len(got))
			}

		})
	}
}

func TestCategoriesService(t *testing.T) {

	services, err := models.NewServices(
		models.WithGorm(
			Dialect(),
			ConnectionInfo()),
		models.WithLogMode(false),
		models.WithCategory(),
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())

	t.Run("Find", testCategoriesService_Find(services.Category, services.GetCategoriesSeed))
}

func testCategoriesService_Find(cs models.CategoryService, getCategoriesSeed func() []models.Category) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("FindAll", func(t *testing.T) {
			want := len(getCategoriesSeed())
			got, err := cs.FindAll()
			if err != nil {
				t.Fatalf("categories could not be fetched error: %s", err.Error())
			}
			if len(got) <= 0 {
				t.Fatalf("expected to find %d categories but got %d categories", want, len(got))
			}

		})
	}
}

func testJobsService_Delete(jobPostService models.JobPostService) func(t *testing.T) {
	return func(t *testing.T) {
		if err := jobPostService.Delete(1); err != nil {
			t.Error(err)
		}
	}
}

func testJobsService_Update(jobPostService models.JobPostService, skillsService models.SkillsService) func(t *testing.T) {
	return func(t *testing.T) {
		got := findJobByID(jobPostService, 1, t)

		t.Run("Fields", func(t *testing.T) {
			got.Title = "Golang Dev Wanted 2"
			got.Description = "Golang Dev Wanted For Backend Development"
			got.ApplyAt = "samysoft@gmail.com"
			got.LocationID = 1
			got.CategoryID = 1
			if err := jobPostService.Update(got); err != nil {
				t.Fatal(err)
			}
			want := findJobByID(jobPostService, 1, t)
			compareJobPostsFields(*got, want, t)
		})

		testAddSkill(t, skillsService, got)

		testRemoveSkill(t, skillsService, got)
		if len(got.Skills) != 0 {
			t.Errorf("expected skills list to be empty, but got = %v elements", len(got.Skills))
		}
	}
}

func testJobsService_Find(jobPostService models.JobPostService) func(t *testing.T) {
	return func(t *testing.T) {
		want := models.JobPost{
			Title:       "Golang Dev Wanted",
			Description: "Golang Dev Wanted For API Development",
			ApplyAt:     "samysoft@gmail.com",
			UserID:      1,
			LocationID:  2,
			CategoryID:  2,
		}

		t.Run("ByID", func(t *testing.T) {
			want.ID = 1
			got := findJobByID(jobPostService, want.ID, t)
			if want.ID != got.ID {
				t.Fatalf("invalid job retrieved, got job with ID:%v, want job with ID %v", got.ID, want.ID)
			}

			compareJobPostsFields(want, got, t)

		})

		t.Run("ByUserID", func(t *testing.T) {

			got := findJobsByUserID(jobPostService, want.UserID, t)
			if len(got) <= 0 {
				t.Fatalf("expected to find %d job posts for user ID %d, but got %d job posts", 1, want.UserID, len(got))
			}
			compareJobPostsFields(want, &got[0], t)

		})
	}
}

func compareJobPostsFields(want models.JobPost, got *models.JobPost, t *testing.T) {
	if want.Title != got.Title {
		t.Errorf("invalid job Title, got job with Title:%v, want job with Title %v", got.Title, want.Title)
	}
	if want.Description != got.Description {
		t.Errorf("invalid job Description, got job with Description:%v, want job with Description %v", got.Description, want.Description)
	}
	if want.ApplyAt != got.ApplyAt {
		t.Errorf("invalid job ApplyAt, got job with ApplyAt:%v, want job with ApplyAt %v", got.ApplyAt, want.ApplyAt)
	}
	if want.LocationID != got.LocationID {
		t.Errorf("invalid job LocationID, got job with LocationID:%v, want job with LocationID %v", got.LocationID, want.LocationID)
	}
	if want.UserID != got.UserID {
		t.Errorf("invalid job UserID, got job with UserID:%v, want job with UserID %v", got.UserID, want.UserID)
	}
	if want.CategoryID != got.CategoryID {
		t.Errorf("invalid job CategoryID, got job with CategoryID:%v, want job with CategoryID %v", got.CategoryID, want.CategoryID)
	}
}

func findJobsByUserID(jobPostService models.JobPostService, id uint, t *testing.T) []models.JobPost {
	got, err := jobPostService.ByUserID(id)
	if err != nil {
		t.Error(err)
	}
	return got
}
func findJobByID(jobPostService models.JobPostService, id uint, t *testing.T) *models.JobPost {
	got, err := jobPostService.ByID(id)
	if err != nil {
		t.Error(err)
	}
	return got
}

type SadPathCreateTest struct {
	name            string
	jobPostModifier func(jp *models.JobPost)
	expectedError   error
}

func testJobsService_Create(jobPostService models.JobPostService) func(t *testing.T) {
	return func(t *testing.T) {
		newJobPost := mockJobPost()
		if err := jobPostService.Create(&newJobPost); err != nil {
			t.Error(err)
		}

		sadPathsTests := []SadPathCreateTest{
			{
				name: "user ID is required",
				jobPostModifier: func(jp *models.JobPost) {
					jp.UserID = 0
				},
				expectedError: models.ErrUserIDRequired,
			},
			{
				name: "Description is required",
				jobPostModifier: func(jp *models.JobPost) {
					jp.Description = ""
				},
				expectedError: models.ErrDescriptionRequired,
			},
			{
				name: "Title is required",
				jobPostModifier: func(jp *models.JobPost) {
					jp.Title = ""
				},
				expectedError: models.ErrTitleRequired,
			},
			{
				name: "ApplyAt is required",
				jobPostModifier: func(jp *models.JobPost) {
					jp.ApplyAt = ""
				},
				expectedError: models.ErrApplyAtRequired,
			},
			{
				name: "LocationID is required",
				jobPostModifier: func(jp *models.JobPost) {
					jp.LocationID = 0
				},
				expectedError: models.ErrLocationIDRequired,
			},
			{
				name: "CategoryID is required",
				jobPostModifier: func(jp *models.JobPost) {
					jp.CategoryID = 0
				},
				expectedError: models.ErrCategoryIDRequired,
			},
		}
		for _, spt := range sadPathsTests {
			t.Run(fmt.Sprintf("SadPath: %s", spt.name), func(t *testing.T) {
				jp := mockJobPost()
				spt.jobPostModifier(&jp)
				wantError := spt.expectedError
				if err := jobPostService.Create(&jp); err != wantError {
					t.Errorf("should return %q error got %q error", wantError, err)
				}
			})
		}

	}
}

func mockJobPost() models.JobPost {
	return models.JobPost{
		Title:       "Golang Dev Wanted",
		Description: "Golang Dev Wanted For API Development",
		ApplyAt:     "samysoft@gmail.com",
		UserID:      1,
		LocationID:  2,
		CategoryID:  2,
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
