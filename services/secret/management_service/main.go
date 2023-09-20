package main

import (
	"encoding/base64"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CreateSchema struct {
	Name     string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ReadSchema struct {
	ID string `json:"id" binding:"required"`
}

type SecretModel struct {
	gorm.Model
  	Name  string
  	Value string
	UUID string
}

func Base64Encode(input string) string {
    data := []byte(input)

    encoded := base64.StdEncoding.EncodeToString(data)

    return encoded
}

func Base64Decode(encoded string) (string, error) {
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return "", err
    }

    decodedStr := string(decoded)

    return decodedStr, nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  	
	if err != nil {
    	panic("failed to connect database")
  	}

  	db.AutoMigrate(&SecretModel{})

    router := gin.Default()

    router.POST("/create", func(context *gin.Context) {
		var body CreateSchema

		if err := context.BindJSON(&body); err != nil {
			context.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
			
			return
		}

		encodedValue := Base64Encode(body.Value)

		newUUID, err := exec.Command("uuidgen").Output()
		
		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
			
			return
		}

		db.Create(&SecretModel{Name: body.Name, Value: encodedValue, UUID: string(newUUID)})

        context.JSON(http.StatusAccepted, gin.H{
            "id": string(newUUID),
        })
    })

	router.POST("/read", func(context *gin.Context) {
		var body ReadSchema

		if err := context.BindJSON(&body); err != nil {
			context.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
			
			return
		}

		var secretModel SecretModel
  		
		db.First(&secretModel, "uuid = ?", body.ID) 

		decodedValue, err := Base64Decode(secretModel.Value)

		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{
					"error": err.Error(),
				},
			)
			
			return
		}

        context.JSON(http.StatusAccepted, gin.H{
            "value": decodedValue,
        })
    })

    router.Run(":8080")
}
