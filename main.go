package main

import (
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/uploadFile/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MyFile struct {
	Id   int                   `uri:"id"`
	Name string                `form:"name"`
	Img  *multipart.FileHeader `form:"img"`
}

func main() {
	DB := ConDB()
	r := gin.Default()
	r.POST("/user/:id", func(c *gin.Context) {
		var myFile MyFile
		var SaveFile models.MyFile
		if err := c.ShouldBind(&myFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if err := c.SaveUploadedFile(myFile.Img, "assets/"+myFile.Img.Filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		SaveFile.Id = myFile.Id
		SaveFile.Name = myFile.Name
		SaveFile.Img = "assets/" + myFile.Img.Filename
		db := DB.Create(SaveFile)
		if db.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": db.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, myFile)
	})

	r.DELETE("/user/:id", func(c *gin.Context) {
		var SaveFile models.MyFile
		// "assets/"+myFile.Img.Filename

		// SaveFile.Img = "assets/" + myFile.Img.Filename
		stringId := c.Param("id")
		dbFile := DB.Find(&SaveFile, stringId)
		if dbFile.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbFile.Error.Error()})
			return
		}

		err := os.Remove(SaveFile.Img)
		log.Println(SaveFile.Img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		db := DB.Delete(SaveFile)
		if db.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": db.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": "Deleted file has been successfully"})
	})

	log.Println("Application runnning on http://localhost:8080")
	r.Run(":8080")
}

// uploadFile

func ConDB() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/uploadFile?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var files models.MyFile
	// Migrate the schema
	db.AutoMigrate(&files)
	log.Println("connection to database successfully")
	return db
}
