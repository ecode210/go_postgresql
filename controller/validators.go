package controller

import (
	"github.com/go-playground/validator/v10"
	"time"
)

// AboveAge - Checks if Date of birth is greater than 16 years ago
var AboveAge validator.Func = func(fl validator.FieldLevel) bool {
	dob, ok := fl.Field().Interface().(int64)
	if ok {
		if dob <= time.Now().Add(-time.Hour * 140160).Unix() {
			return true
		}
	}
	return false
}