package model

import (
	"regexp"

	"github.com/pkg/errors"
)

var productNameRegex = regexp.MustCompile("^[a-zA-Z0-9 ]{3,50}")

type ProductName string

func (p ProductName) IsValid() error {
	if !productNameRegex.MatchString(string(p)) {
		return errors.Errorf("invalid product name")
	}
	return nil
}

type Product struct {
	PID   int64  `json:"pid"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}
