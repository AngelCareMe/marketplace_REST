package entity

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID             uuid.UUID `json:"id"`
	Header         string    `json:"header"`
	Content        string    `json:"content"`
	Image          string    `json:"image"`
	Price          float64   `json:"price"`
	AuthorID       uuid.UUID `json:"author_id"`
	CreatedAt      time.Time `json:"created_at"`
	IsOwnPost      bool      `json:"is_own_post"`
	AuthorUsername string    `json:"author_username"`
}

func (p *Post) Validate() error {
	if p.Header == "" {
		return fmt.Errorf("header can't be empty")
	}
	if len(strings.TrimSpace(p.Header)) < 5 {
		return fmt.Errorf("header must be at least 5 characters long")
	}
	if len(p.Header) > 100 {
		return fmt.Errorf("header must not exceed 100 characters")
	}

	if p.Content == "" {
		return fmt.Errorf("content can't be empty")
	}
	if len(strings.TrimSpace(p.Content)) < 10 {
		return fmt.Errorf("content must be at least 10 characters long")
	}
	if len(p.Content) > 1000 {
		return fmt.Errorf("content must not exceed 1000 characters")
	}

	if p.Image == "" {
		return fmt.Errorf("post must have image")
	}
	if _, err := url.ParseRequestURI(p.Image); err != nil {
		return fmt.Errorf("image must be a valid URL: %w", err)
	}
	validImageExt := regexp.MustCompile(`\.(png|jpg|jpeg)$`)
	if !validImageExt.MatchString(strings.ToLower(p.Image)) {
		return fmt.Errorf("image must be in PNG or JPEG format")
	}

	if p.Price < 0 {
		return fmt.Errorf("price must be positive")
	}
	if p.Price > 1000000 {
		return fmt.Errorf("price must not exceed 1000000")
	}
	return nil
}
