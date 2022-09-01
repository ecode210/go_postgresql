package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_postgresql/config"
	"go_postgresql/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

var Validate *validator.Validate

func Home(c *gin.Context){
	c.String(http.StatusOK, "Welcome to Go Postgresql")
}

func GetAllUsers(c *gin.Context){
	var allUsers []gin.H
	var users []model.User
	err := config.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err,
		})
		return
	}

	for i := range users{
		allUsers = append(allUsers, gin.H{
			"full_name": users[i].FullName,
			"email": users[i].Email,
			"phone_number": users[i].PhoneNumber,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error":  false,
		"data":   allUsers,
	})
}

func RegisterUser(c *gin.Context){
	var input model.User

	// Binds & Validates input from request body to User struct
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := ""
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "FullName":
				errorMessage = errorMessage + "Full Name needs to be greater than 5 characters and less than 30 characters. "
				break
			case "DateOfBirth":
				errorMessage = errorMessage + "User is below 16 years. "
				break
			case "PhoneNumber":
				errorMessage = errorMessage + "Invalid Phone Number (Try adding country code e.g '+234'). "
				break
			case "Email":
				errorMessage = errorMessage + "Invalid Email address. "
				break
			case "Password":
				errorMessage = errorMessage + "Password too short/long. "
				break
			}
		}
		if errorMessage == "" {
			errorMessage = err.(validator.ValidationErrors).Error()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": strings.Trim(errorMessage, " "),
		})
		return
	}

	// Checks if a User with the same email already exists
	var existingUser model.User
	err = config.DB.First(&existingUser, "email = ?", input.Email).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}
	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "User already exists",
		})
		return
	}

	// Hashes Password
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}
	input.Password = string(password)

	// Adds New User to Postgresql DB
	err = config.DB.Create(&input).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		fmt.Println("Error:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"message": "Registration successful!",
	})
}

func LoginUser(c *gin.Context){
	// Binds & Validates input from request body to UserLogin struct
	var input model.UserLogin
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := ""
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Email":
				errorMessage = errorMessage + "Invalid Email address. "
				break
			case "Password":
				errorMessage = errorMessage + "Invalid Password. "
				break
			}
		}
		if errorMessage == "" {
			errorMessage = err.(validator.ValidationErrors).Error()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": strings.Trim(errorMessage, " "),
		})
		return
	}

	var allUsers []model.User
	err = config.DB.Find(&allUsers).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	// Checks if local DB is empty
	if len(allUsers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "No Users available",
		})
		return
	}

	// Verifies Login credentials
	for i := range allUsers {
		if allUsers[i].Email == input.Email {
			err = bcrypt.CompareHashAndPassword([]byte(allUsers[i].Password), []byte(input.Password))
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"error": false,
					"message": "Logged in!",
					"data": gin.H{
						"id": allUsers[i].ID,
					},
				})
				return
			}else{
				c.JSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error": true,
					"message": "Incorrect Password",
				})
				return
			}
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"error": true,
		"message": "Wrong Email/Password",
	})
}

func UpdateDetails(c *gin.Context){
	// Binds input from request body to UserUpdate struct
	var input model.UserUpdate
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	// Finds user in database
	var user model.User
	err = config.DB.First(&user, "id = ?", c.Param("id")).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "User not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	// Checks if user details were provided
	if input.FullName == "" && input.PhoneNumber == "" && input.DateOfBirth == 0 {
		c.JSON(http.StatusAccepted, gin.H{
			"status": http.StatusAccepted,
			"error": false,
			"message": "No User details updated!",
		})
		return
	}

	// Validates Full Name
	if input.FullName != "" {
		err = Validate.Var(input.FullName, "gte=5,lte=30")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"error":   true,
				"message": "Full Name needs to be greater than 5 characters and less than 30 characters.",
			})
			return
		}
	}

	// Validates Phone Number
	if input.PhoneNumber != "" {
		err = Validate.Var(input.PhoneNumber, "e164")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"error":   true,
				"message": "Invalid Phone Number (Try adding country code e.g '+234').",
			})
			return
		}
	}

	// Validates Date of birth
	if input.DateOfBirth != 0 {
		// Checks if Date of birth is greater than 16 years ago
		timeValidator := fmt.Sprintf("lte=%v", time.Now().Add(-time.Hour*140160).Unix())
		err = Validate.Var(input.DateOfBirth, timeValidator)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"error":   true,
				"message": "User is below 16 years.",
			})
			return
		}
	}

	// Saves new updates to database
	user.UpdateUser(input)
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"message": "User details updated!",
	})
}

func DeleteUser(c *gin.Context) {
	// Finds index of user in local DB from :id in endpoint URL
	var user model.User
	err := config.DB.First(&user, "id = ?", c.Param("id")).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "User not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	err = config.DB.Delete(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"message": "User deleted! (We're keeping your data though😈)",
	})
}

