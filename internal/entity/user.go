package entity

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username can't be empty")
	}
	if len(strings.TrimSpace(u.Username)) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(u.Username) > 50 {
		return fmt.Errorf("username must not exceed 50 characters")
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(u.Username) {
		return fmt.Errorf("username can only contain letters, digits, and underscores")
	}

	return nil
}
