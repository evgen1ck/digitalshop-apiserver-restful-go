package validator

import (
	"errors"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

func Validate(s string, validators ...func(string) error) error {
	for _, v := range validators {
		if err := v(s); err != nil {
			return err
		}
	}
	return nil
}

func IsMinMaxLen(min, max int) func(string) error {
	return func(str string) error {
		l := len(str)
		if l < min {
			return errors.New("the value is too short (minimum is " + string(rune(min)) + " characters)")
		} else if l > max {
			return errors.New("the value is too long (maximum is " + string(rune(max)) + " characters)")
		}
		return nil
	}
}

func IsEmail() func(string) error {
	return func(str string) error {
		_, err := mail.ParseAddress(str)
		if err != nil {
			return errors.New("the value is not an email")
		}
		return nil
	}
}

func IsUrl() func(string) error {
	return func(str string) error {
		_, err := url.ParseRequestURI(str)
		if err != nil {
			return errors.New("the value is not a url")
		}
		return nil
	}
}

func IsInt64() func(string) error {
	return func(str string) error {
		_, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return errors.New("the value is not an integer data type")
		}
		return nil
	}
}

func IsUint64() func(string) error {
	return func(str string) error {
		_, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return errors.New("the value is not an unsigned integer data type")
		}
		return nil
	}
}

func IsBlank() func(string) error {
	return func(str string) error {
		if strings.TrimSpace(str) == "" {
			return errors.New("the value is blank")
		}
		return nil
	}
}

func IsContainsSpaces() func(string) error {
	return func(str string) error {
		for _, c := range str {
			if unicode.IsSpace(c) {
				return errors.New("the value contains a space(s)")
			}
		}
		return nil
	}
}
