package controllers

import (
	"database/sql"
	"final-project/config"
	"final-project/structs"
	"final-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Register a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body structs.User true "User Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/register [post]
func Register(c *gin.Context) {
	var user structs.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	query := `INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4)`
	_, err = config.DB.Exec(query, user.Name, user.Email, string(hash), "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body structs.User true "User login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/login [post]
func Login(c *gin.Context) {
	var creds structs.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user structs.User
	query := `SELECT id, name, email, password_hash, role FROM users WHERE email = $1`
	row := config.DB.QueryRow(query, creds.Email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
