package controllers

import (
	"final-project/config"
	"final-project/structs"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBooks(c *gin.Context) {
	rows, err := config.DB.Query(`SELECT id, title, author, stock, category_id FROM books`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch books"})
		return
	}
	defer rows.Close()

	var books []structs.Book
	for rows.Next() {
		var b structs.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Stock, &b.CategoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error scanning books"})
			return
		}
		books = append(books, b)
	}

	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var book structs.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = uuid.New().String()
	fmt.Println("New book ID:", book.ID)

	query := `INSERT INTO books (id, title, author, stock, category_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := config.DB.Exec(query, book.ID, book.Title, book.Author, book.Stock, book.CategoryID)
	if err != nil {
		fmt.Println("Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert book"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book created", "id": book.ID})
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book structs.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE books SET title=$1, author=$2, stock=$3, category_id=$4 WHERE id=$5`
	res, err := config.DB.Exec(query, book.Title, book.Author, book.Stock, book.CategoryID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book updated"})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	res, err := config.DB.Exec(`DELETE FROM books WHERE id=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}
