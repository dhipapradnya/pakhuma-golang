package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Age       string    `gorm:"column:age"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func main() {
	dsn := "root:@tcp(localhost:3306)/open_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database")
	}

	router := gin.Default()

	router.GET("/users", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		db.Create(&user)
		c.JSON(http.StatusCreated, gin.H{"data": user})
	})

	router.GET("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.DELETE("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		db.Delete(&user)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	})

	router.PUT("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		result := db.First(&user, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		user.UpdatedAt = time.Now()
		db.Save(&user)
		c.JSON(http.StatusOK, gin.H{"data": user})
	})

	router.Run(":3000")
}
