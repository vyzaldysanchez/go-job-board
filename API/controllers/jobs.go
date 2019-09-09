package controllers

import (
	"net/http"

	"github.com/samueldaviddelacruz/go-job-board/API/models"
)

type Jobs struct {
	JobService models.JobPostService
}

func NewJobs(js models.JobPostService) *Jobs {
	return &Jobs{
		js,
	}
}

// GET /jobs
func (j *Jobs) List(w http.ResponseWriter, r *http.Request) {
	jobs, err := j.JobService.FindAll()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, jobs)
}

//POST /Create
func (j *Jobs) Create(w http.ResponseWriter, r *http.Request) {

	jobPost := models.JobPost{

	}
	parseJSON(w, r, &jobPost)

	if err := j.JobService.Create(&jobPost); err != nil {

		respondJSON(w, http.StatusInternalServerError, "Could not create jobPost")
		return
	}
	respondJSON(w, http.StatusCreated, jobPost)
}

//POST /Update
func (j *Jobs) Update(w http.ResponseWriter, r *http.Request) {

	jobPost := models.JobPost{

	}
	parseJSON(w, r, &jobPost)

	if err := j.JobService.Update(&jobPost); err != nil {

		respondJSON(w, http.StatusInternalServerError, "Could not update jobPost")
		return
	}
	respondJSON(w, http.StatusCreated, jobPost)
}
