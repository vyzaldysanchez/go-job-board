package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

type Jobs struct {
	js models.JobPostService
	ss models.SkillsService
}

func NewJobs(js models.JobPostService, ss models.SkillsService) *Jobs {
	return &Jobs{
		js,
		ss,
	}
}

// GET /jobs
func (j *Jobs) List(w http.ResponseWriter, r *http.Request) {
	queryObj := models.JobPost{}
	queryObj.Title = r.URL.Query().Get("q")
	if userId, err := strconv.Atoi(r.URL.Query().Get("u")); err == nil {
		queryObj.UserID = uint(userId)
	}
	if locationId, err := strconv.Atoi(r.URL.Query().Get("l")); err == nil {
		queryObj.LocationID = uint(locationId)
	}
	if categoryId, err := strconv.Atoi(r.URL.Query().Get("c")); err == nil {
		queryObj.CategoryID = uint(categoryId)
	}
	queryObj.Skills = extractSkillsFromQueryStr(r)

	jobs, err := j.js.FindAll(queryObj)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, jobs)
}

func extractSkillsFromQueryStr(r *http.Request) []models.Skill {
	var skills []models.Skill
	if skillsStr := r.URL.Query().Get("sk"); skillsStr != "" {
		skillsIds := strings.Split(skillsStr, ",")
		for _, skillId := range skillsIds {
			if id, err := strconv.Atoi(skillId); err == nil {
				newSkill := models.Skill{}
				newSkill.ID = uint(id)
				skills = append(skills, newSkill)
			}
		}
	}
	return skills
}

//POST /jobs
func (j *Jobs) Create(w http.ResponseWriter, r *http.Request) {

	jobPost := models.JobPost{

	}
	err := parseJSON(r, &jobPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := j.js.Create(&jobPost); err != nil {

		respondJSON(w, http.StatusInternalServerError, "Could not create jobPost")
		return
	}
	respondJSON(w, http.StatusCreated, jobPost)
}

//PUT /jobs/id
func (j *Jobs) Update(w http.ResponseWriter, r *http.Request) {

	jobPost, err := j.getJobByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	parseJSON(r, jobPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := j.js.Update(jobPost); err != nil {

		respondJSON(w, http.StatusInternalServerError, "Could not update jobPost")
		return
	}
	respondJSON(w, http.StatusCreated, jobPost)
}

//DELETE /jobs/id
func (j *Jobs) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, "Could not delete jobPost")
		return
	}
	if err := j.js.Delete(uint(id)); err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, fmt.Sprintf("Removed Jobpost with ID %v", id))
}

// PUT /jobs/id/add-skill
func (j *Jobs) AddJobPostSkill(w http.ResponseWriter, r *http.Request) {

	jobPost, err := j.getJobByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	skill := models.Skill{}
	parseJSON(r, &skill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := j.ss.AddSkillToOwner(jobPost, skill); err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, "skills updated successfully")
}

// PUT /user/id/remove-skill
func (j *Jobs) RemoveJobPostSkill(w http.ResponseWriter, r *http.Request) {
	jobPost, err := j.getJobByID(r)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	skill := models.Skill{}
	parseJSON(r, &skill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := j.ss.DeleteSkillFromOwner(jobPost, skill); err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, "skills updated successfully")
}

func (j *Jobs) getJobByID(r *http.Request) (*models.JobPost, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {

		return nil, err
	}
	jobPost, err := j.js.ByID(uint(id))
	if err != nil {

		return nil, err
	}
	return jobPost, nil
}
