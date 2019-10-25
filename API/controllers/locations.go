package controllers

import (
	"github.com/samueldaviddelacruz/go-job-board/API/models"
	"net/http"
)

type Locations struct {
	ls models.LocationService
}

func NewLocations(ls models.LocationService) *Locations {
	return &Locations{
		ls,
	}
}

// GET /locations
func (c *Locations) List(w http.ResponseWriter, r *http.Request) {
	locations, err := c.ls.FindAll()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, locations)
}
