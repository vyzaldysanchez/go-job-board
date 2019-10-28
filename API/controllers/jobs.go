package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

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
	testJ := models.JobPost{
		UserID:      0,
		Title:       r.URL.Query().Get("q"),
		Location:    nil,
		LocationID:  0,
		Category:    nil,
		CategoryID:  0,
		Description: "",
		ApplyAt:     "",
		Skills:      nil,
	}
	if locationId, err := strconv.Atoi(r.URL.Query().Get("l")); err == nil {
		testJ.LocationID = uint(locationId)
	}
	//jobs, err := j.js.FindAll()
	jobs, err := j.js.FindByParameters(&testJ)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, jobs)
}

//POST /jobs
func (j *Jobs) Create(w http.ResponseWriter, r *http.Request) {

	jobPost := models.JobPost{

	}
	parseJSON(w, r, &jobPost)

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

	parseJSON(w, r, jobPost)

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
	parseJSON(w, r, &skill)

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
	parseJSON(w, r, &skill)

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
