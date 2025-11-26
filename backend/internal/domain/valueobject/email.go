package valueobject

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if !emailRegex.MatchString(email) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: email}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Value() string {
	return e.value
}



