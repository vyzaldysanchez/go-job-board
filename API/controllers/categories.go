package controllers

import (
	"github.com/samueldaviddelacruz/go-job-board/API/models"
	"net/http"
)

type Categories struct {
	cs models.CategoryService
}

func NewCategories(cs models.CategoryService) *Categories {
	return &Categories{
		cs,
	}
}

// GET /categories
func (c *Categories) List(w http.ResponseWriter, r *http.Request) {
	categories, err := c.cs.FindAll()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, categories)
}
