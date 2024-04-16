// Package validator provides utilities for validating data.
package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator is a struct that holds field errors.
type Validator struct {
	FieldErrors    map[string]string // FieldErrors is a map of field names to error messages.
	NonFieldErrors []string
}

// Valid checks if the validator has any field errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError adds an error message for a field to the validator.
func (v *Validator) AddFieldError(key, message string) {

	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField checks a condition and, if it's not met, adds an error message for a field to the validator.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank checks if a string is not blank.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxRunes checks if a string has no more than a certain number of runes.
func MaxRunes(value string, maxCount int) bool {
	return utf8.RuneCountInString(value) <= maxCount
}

// AllowedInt checks if an integer is in a list of allowed values.
func AllowedInt(value int, allowedValues ...int) bool {
	for _, allowedVal := range allowedValues {
		if allowedVal == value {
			return true
		}
	}
	return false
}

func MinRunes(value string, minCount int) bool {
	return utf8.RuneCountInString(value) >= minCount
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
