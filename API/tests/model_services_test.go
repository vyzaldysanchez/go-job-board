package model_services_test

import (
	"fmt"
	"os"
	"reflect"
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
	)
	must(err)

	defer services.Close()
	must(services.DestructiveReset())
	t.Run("Create", testUserService_Create(services.User))
	t.Run("Find", testUserService_Find(services.User))
	t.Run("Update", testUserService_Update(services.User))
	t.Run("Delete", testUserService_Delete(services.User))

	// teardown
}

func testUserService_Update(us models.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		want := findUserByID(us, 1, t)

		t.Run("CompanyProfile", func(t *testing.T) {
			got := findUserByID(us, 1, t)
			got.CompanyProfile = &models.CompanyProfile{
				Website:        "samysoft.com",
				FoundedYear:    1991,
				Description:    "a very nice company",
				CompanyLogoUrl: "a logo url",
			}

			if err := us.Update(got); err != nil {
				t.Error(err)
			}

			want.CompanyProfile = &models.CompanyProfile{
				Website:        "samysoft.com",
				FoundedYear:    1991,
				Description:    "a very nice company",
				CompanyLogoUrl: "a logo url",
			}

			if reflect.DeepEqual(want.CompanyProfile, got.CompanyProfile) {
				t.Error("Website did not update correctly")
			}

		})

	}
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
