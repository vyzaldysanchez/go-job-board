package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup
func NewUsers(us models.UserService, ss models.SkillsService) *Users {
	return &Users{
		us: us,
		ss: ss,
	}
}

// Users Represents a Users controller
type Users struct {
	us models.UserService
	ss models.SkillsService
}
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PUT /user/id
func (u *Users) Update(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, "Invalid user ID")
		return
	}

	err = parseJSON(r, companyUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := u.us.Update(companyUser); err != nil {
		//vd.SetAlert(err)
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, "resource updated successfully")
}

// PUT /user/id/company-profile
func (u *Users) UpdateCompanyProfile(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, "Invalid user ID")
		return
	}
	newCompanyProfile := &models.CompanyProfile{}
	err = parseJSON(r, newCompanyProfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if companyUser.CompanyProfile == nil {
		companyUser.CompanyProfile = &models.CompanyProfile{}
	}
	companyUser.CompanyProfile.Description = newCompanyProfile.Description
	companyUser.CompanyProfile.Website = newCompanyProfile.Website
	companyUser.CompanyProfile.CompanyLogoUrl = newCompanyProfile.CompanyLogoUrl
	companyUser.CompanyProfile.FoundedYear = newCompanyProfile.FoundedYear

	if err := u.us.Update(companyUser); err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, "resource updated successfully")
}

// PUT /user/id/company-profile/add-skill
func (u *Users) AddCompanyProfileSkill(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	skill := models.Skill{}
	err = parseJSON(r, &skill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if companyUser.CompanyProfile != nil {
		if err := u.ss.AddSkillToOwner(companyUser.CompanyProfile, skill); err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	respondJSON(w, http.StatusCreated, "skills updated successfully")
}

// PUT /user/id/company-profile/remove-skill
func (u *Users) RemoveCompanyProfileSkill(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, "Invalid user ID")
		return
	}
	skill := models.Skill{}
	err = parseJSON(r, &skill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if companyUser.CompanyProfile != nil {
		if err := u.ss.DeleteSkillFromOwner(companyUser.CompanyProfile, skill); err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondJSON(w, http.StatusCreated, "skills updated successfully")
}

// PUT /user/id/company-profile/add-benefit
func (u *Users) AddCompanyProfileBenefit(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	benefit := models.CompanyBenefit{}
	err = parseJSON(r, &benefit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if companyUser.CompanyProfile != nil {
		if err := u.us.AddCompanyProfileBenefit(companyUser.CompanyProfile, benefit); err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondJSON(w, http.StatusCreated, "benefit added successfully")
}

// PUT /user/id/company-profile/remove-benefit
func (u *Users) RemoveCompanyProfileBenefit(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	benefit := models.CompanyBenefit{}
	err = parseJSON(r, &benefit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if companyUser.CompanyProfile != nil {
		if err := u.us.RemoveCompanyProfileBenefit(companyUser.CompanyProfile, benefit); err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondJSON(w, http.StatusCreated, "benefit removed successfully")
}

// PUT /user/id/company-profile/update-benefit
func (u *Users) UpdateCompanyProfileBenefit(w http.ResponseWriter, r *http.Request) {
	companyUser, err := u.getUserByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	benefit := &models.CompanyBenefit{}
	err = parseJSON(r, benefit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if companyUser.CompanyProfile != nil {
		benefit.CompanyProfileID = companyUser.CompanyProfile.ID
		if err := u.us.UpdateCompanyProfileBenefit(benefit); err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondJSON(w, http.StatusCreated, "benefit updated successfully")
}

func (u *Users) getUserByID(r *http.Request) (*models.User, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {

		return nil, err
	}
	companyUser, err := u.us.ByID(uint(id))
	if err != nil {

		return nil, err
	}
	return companyUser, nil
}
