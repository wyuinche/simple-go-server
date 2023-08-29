package model

import (
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	passwordEntireRegex     = regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*+]{8,}$")
	passwordSpecialChrRegex = regexp.MustCompile("[!@#$%^&*+]+")
	passwordNumRegex        = regexp.MustCompile("[0-9]+")
	passwordChrRegex        = regexp.MustCompile("[a-zA-Z]+")
)

type Password string

func (p Password) IsValid() error {
	e := errors.Errorf("invalid password")

	pw := string(p)

	if !passwordEntireRegex.MatchString(pw) {
		return errors.WithMessage(e, "invalid character or short length")
	}

	if !passwordSpecialChrRegex.MatchString(pw) {
		return errors.WithMessage(e, "no special character")
	}

	if !passwordNumRegex.MatchString(pw) {
		return errors.WithMessage(e, "no digit")
	}

	if !passwordChrRegex.MatchString(pw) {
		return errors.WithMessage(e, "no alphabet character")
	}

	return nil
}

func (p Password) Hash() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(bytes), err
}

func (p Password) CompareWithHash(hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
	if err != nil {
		return false
	} else {
		return true
	}
}

var userIDRegex = regexp.MustCompile("^[a-zA-Z0-9]{3,18}$")

type UserID string

func (id UserID) IsValid() error {
	if !userIDRegex.MatchString(string(id)) {
		return errors.Errorf("invalid user id")
	}
	return nil
}

const (
	RoleUser    = "user"
	RoleManager = "manager"
)

type User struct {
	UID      int64  `json:"uid"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
