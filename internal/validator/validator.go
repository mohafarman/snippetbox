package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

/* This pattern is parsed once at startup and stored to be used */
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	/* If the key already exists, we don't add a new error message */
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	/* Add error message to field error only if the validation check is not 'ok' */
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	/* Return true if value is not empty string */
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	/* Return true if value contains no more than n characters */
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	/* Return true if value contains less than n characters */
	return utf8.RuneCountInString(value) >= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
