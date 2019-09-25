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
		t.Error(err)
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
	t.Run("add-skill", func(t *testing.T) {
		skill := models.Skill{}
		skill.ID = 1

		if err := ss.AddSkillToOwner(got, skill); err != nil {
			t.Errorf("error adding company profile skill error = %v, companyProfile = %v, skill = %v", err, got, skill)
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
	t.Run("remove-skill", func(t *testing.T) {
		skill := models.Skill{}
		skill.ID = 1

		if err := ss.DeleteSkillFromOwner(got, skill); err != nil {
			t.Errorf("error removing company profile skill error = %v, companyProfile = %v, skill = %v", err, got, skill)
		}
		if len(got.Skills) != 0 {
			t.Errorf("expected skills list to be empty, but got = %v elements", len(got.Skills))
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
