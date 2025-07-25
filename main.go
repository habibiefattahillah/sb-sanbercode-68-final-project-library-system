package main

import (
	"final-project/config"
	"final-project/controllers"
	"final-project/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/users/register", controllers.Register)
		api.POST("/users/login", controllers.Login)
		api.GET("/books", controllers.GetBooks)
	}

	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetString("user_id"),
				"email":   c.GetString("email"),
				"role":    c.GetString("role"),
			})
		})
		
		protected.POST("/books", controllers.CreateBook)
		protected.PUT("/books/:id", controllers.UpdateBook)
		protected.DELETE("/books/:id", controllers.DeleteBook)

		protected.GET("/categories", controllers.GetCategories)
		protected.POST("/categories", controllers.CreateCategory)
	
		protected.POST("/books/:book_id/borrow", controllers.BorrowBook)
		protected.POST("/books/:book_id/return", controllers.ReturnBook)
		protected.GET("/borrowings", controllers.GetBorrowings)
	}

	r.Run(":" + config.EnvPort())
}
