package controllers

import (
	"database/sql"
	"final-project/config"
	"final-project/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// BorrowBook godoc
// @Summary Borrow a book
// @Description Allows a user to borrow a book if stock is available
// @Tags borrowings
// @Produce json
// @Param book_id path string true "Book ID to borrow"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /books/{book_id}/borrow [post]
func BorrowBook(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("book_id")

	var stock int
	err := config.DB.QueryRow(`SELECT stock FROM books WHERE id = $1`, bookID).Scan(&stock)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	if stock < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book is out of stock"})
		return
	}

	var exists bool
	err = config.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM borrowings
			WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
		)
	`, userID, bookID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check borrowings"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you have already borrowed this book"})
		return
	}

	_, err = config.DB.Exec(`
		INSERT INTO borrowings (user_id, book_id, borrowed_at)
		VALUES ($1, $2, $3)
	`, userID, bookID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to borrow book"})
		return
	}

	_, err = config.DB.Exec(`UPDATE books SET stock = stock - 1 WHERE id = $1`, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stock"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book borrowed successfully"})
}

// ReturnBook godoc
// @Summary Return a borrowed book
// @Description Allows a user to return a previously borrowed book
// @Tags borrowings
// @Produce json
// @Param book_id path string true "Book ID to return"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /books/{book_id}/return [post]
func ReturnBook(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("book_id")

	var borrowingID string
	err := config.DB.QueryRow(`
		SELECT id FROM borrowings
		WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL
	`, userID, bookID).Scan(&borrowingID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no active borrow found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query borrowings"})
		}
		return
	}

	_, err = config.DB.Exec(`
		UPDATE borrowings SET returned_at = $1 WHERE id = $2
	`, time.Now(), borrowingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to return book"})
		return
	}

	_, err = config.DB.Exec(`UPDATE books SET stock = stock + 1 WHERE id = $1`, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book returned successfully"})
}

func GetBorrowings(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := config.DB.Query(`
		SELECT b.id, bk.title, b.borrowed_at, b.returned_at
		FROM borrowings b
		JOIN books bk ON b.book_id = bk.id
		WHERE b.user_id = $1
		ORDER BY b.borrowed_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch borrowings"})
		return
	}
	defer rows.Close()

	var result []structs.BorrowingInfo
	for rows.Next() {
		var r structs.BorrowingInfo
		err := rows.Scan(&r.ID, &r.Title, &r.BorrowedAt, &r.ReturnedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning borrowings"})
			return
		}
		result = append(result, r)
	}

	c.JSON(http.StatusOK, result)
}
