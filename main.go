package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetRecord(c *gin.Context, db *gorm.DB) {
	var rec Company
	db.Find(&rec)
	c.JSON(200, rec)

}

func GetRecordById(c *gin.Context, db *gorm.DB) {
	var data Company
	todoID := c.Param("id")

	result := db.First(&data, todoID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(200, data)

}

func AddRecord(c *gin.Context, db *gorm.DB) {
	var data Company
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data"})
		return
	}

	db.Create(&data)

	c.JSON(200, data)
}

func UpdateRecord(c *gin.Context, db *gorm.DB) {
	var data Company
	todoID := c.Param("id")

	result := db.First(&data, todoID)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	var updatedTodo Company
	if err := c.ShouldBindJSON(&updatedTodo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON data"})
		return
	}

	data.Title = updatedTodo.Title
	data.Description = updatedTodo.Description
	db.Save(&data)

	c.JSON(200, data)
}

func DeleteRecord(c *gin.Context, db *gorm.DB) {
	var data Company
	compId := c.Param("id")

	result := db.First(&data, compId)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Todo not found"})
		return
	}

	db.Delete(&data)

	c.JSON(200, gin.H{"message": fmt.Sprintf("Todo with ID %s deleted", compId)})
}

func main() {
	router := gin.Default()

	db, err := gorm.Open(sqlite.Open("Learning.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Company{})
	if err != nil {
		return
	}
	var todos []Company
	fmt.Println(todos)

	router.POST("/comp", func(context *gin.Context) {
		AddRecord(context, db)
	})

	router.GET("/comp", func(context *gin.Context) {
		GetRecord(context, db)
	})

	router.GET("/comp/:id", func(context *gin.Context) {
		GetRecordById(context, db)
	})

	router.PUT("/comp/:id", func(context *gin.Context) {
		UpdateRecord(context, db)
	})

	router.DELETE("/comp/:id", func(context *gin.Context) {
		DeleteRecord(context, db)
	})

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
