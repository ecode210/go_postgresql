package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName    string `json:"full_name,omitempty" binding:"required,gte=5,lte=30"`
	DateOfBirth int64  `json:"date_of_birth,omitempty" binding:"required,above_age"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required,e164"`
	Email       string `json:"email,omitempty" binding:"required,email"`
	Password    string `json:"password,omitempty" binding:"required,gte=8,lte=30"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8,lte=30"`
}

type UserUpdate struct {
	FullName    string `json:"full_name"`
	DateOfBirth int64  `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateUser - Checks for non-empty values from the given UserUpdate and assigns them to the User
func (user *User) UpdateUser(userUpdate UserUpdate) {
	if userUpdate.FullName != "" {
		user.FullName = userUpdate.FullName
	}
	if userUpdate.DateOfBirth != 0 {
		user.DateOfBirth = userUpdate.DateOfBirth
	}
	if userUpdate.PhoneNumber != "" {
		user.PhoneNumber = userUpdate.PhoneNumber
	}
}
