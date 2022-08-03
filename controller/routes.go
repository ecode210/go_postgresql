package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go_postgresql/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

var Validate *validator.Validate

func GetAllUsers(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"data": model.AllUsers,
	})
}

func RegisterUser(c *gin.Context){
	// Binds input from request body to User struct
	var input model.User
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	// Generates UUID
	input.ID = uuid.New().String()

	// Validates User struct from validate tags
	err = Validate.Struct(input)
	if err != nil {
		errorMessage := ""
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "ID":
				errorMessage = errorMessage + "UUID error! Please Try again later. "
				break
			case "FullName":
				errorMessage = errorMessage + "Full Name needs to be greater than 5 characters and less than 30 characters. "
				break
			case "PhoneNumber":
				errorMessage = errorMessage + "Invalid Phone Number (Try adding country code e.g '+234'). "
				break
			case "Email":
				errorMessage = errorMessage + "Invalid Email address. "
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

	// Checks if Date of birth is greater than 16 years ago
	timeValidator := fmt.Sprintf("lte=%v", time.Now().Add(-time.Hour * 140160).Unix())
	err = Validate.Var(input.DateOfBirth, timeValidator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "User is below 16 years.",
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

	// Adds New User to Local DB
	model.AllUsers = append(model.AllUsers, input)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"message": "Registration successful!",
	})
}

func LoginUser(c *gin.Context){
	// Binds input from request body to UserLogin struct
	var input model.UserLogin
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": err.Error(),
		})
		return
	}

	// Checks if local DB is empty
	if model.AllUsers == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "No Users available",
		})
		return
	}

	// Verifies Login credentials
	for i := range model.AllUsers {
		if model.AllUsers[i].Email == input.Email {
			err = bcrypt.CompareHashAndPassword([]byte(model.AllUsers[i].Password), []byte(input.Password))
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"status": http.StatusOK,
					"error": false,
					"message": "Logged in!",
					"data": gin.H{
						"id": model.AllUsers[i].ID,
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

	// Finds index of user in local DB from :id in endpoint URL
	userIndex := -1
	for i := range model.AllUsers{
		if model.AllUsers[i].ID == c.Param("id") {
			userIndex = i
			break
		}
	}
	if userIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "User not found",
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

	model.AllUsers[userIndex].UpdateUser(input)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"message": "User details updated!",
	})
}

func DeleteUser(c *gin.Context) {
	// Finds index of user in local DB from :id in endpoint URL
	userIndex := -1
	for i := range model.AllUsers{
		if model.AllUsers[i].ID == c.Param("id") {
			userIndex = i
			break
		}
	}
	if userIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": true,
			"message": "User not found",
		})
		return
	}

	// Copies the Users below the current User and pastes them into the local DB starting from the current User index
	copy(model.AllUsers[userIndex:], model.AllUsers[userIndex+1:])
	// Sets the last User in the local DB to empty because it is a duplicate of the second-to-the-last User
	model.AllUsers[len(model.AllUsers)-1] = model.User{}
	// Truncates the local DB to remove empty entries
	model.AllUsers = model.AllUsers[:len(model.AllUsers)-1]

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"error": false,
		"message": "User deleted! (We're keeping your data though😈)",
	})
}

