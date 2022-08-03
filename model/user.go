package model

type User struct {
	ID          string `json:"id" validate:"required,uuid"`
	FullName    string `json:"full_name" binding:"required" validate:"gte=5,lte=30"`
	DateOfBirth int64  `json:"date_of_birth" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"  validate:"e164"`
	Email       string `json:"email" binding:"required" validate:"email"`
	Password    string `json:"password" binding:"required" validate:"gte=8,lte=30"`
}

type UserLogin struct {
	Email       string `json:"email" binding:"required" validate:"email"`
	Password    string `json:"password" binding:"required" validate:"gte=8,lte=30"`
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
