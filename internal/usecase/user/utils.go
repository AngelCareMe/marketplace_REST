package usecase

import (
	"fmt"
	"regexp"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(password) > 100 {
		return fmt.Errorf("password must not exceed 100 characters")
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasLetter || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must contain at least one letter, one digit, and one special character")
	}
	return nil
}
