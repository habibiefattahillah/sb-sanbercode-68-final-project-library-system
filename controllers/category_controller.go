package controllers

import (
	"final-project/config"
	"final-project/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCategories godoc
// @Summary Get all categories
// @Description Returns a list of all book categories
// @Tags categories
// @Produce json
// @Success 200 {array} structs.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, name FROM categories`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}
	defer rows.Close()

	var categories []structs.Category
	for rows.Next() {
		var cat structs.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Adds a new category to the system
// @Tags categories
// @Accept json
// @Produce json
// @Param category body structs.Category true "Category to create"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var category structs.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = uuid.New().String()

	_, err := config.DB.Exec(`INSERT INTO categories (id, name) VALUES ($1, $2)`, category.ID, category.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "category created", "id": category.ID})
}
