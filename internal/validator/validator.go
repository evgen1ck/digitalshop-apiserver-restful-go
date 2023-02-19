package validator

import (
	"errors"
)

func Validate(s string, validators ...func(string) error) error {
	for _, v := range validators {
		if err := v(s); err != nil {
			return err
		}
	}
	return nil
}

func MinMaxLenValidator(min, max int) func(string) error {
	return func(s string) error {
		l := len(s)
		if l < min {
			return errors.New("value is too short (minimum is " + string(rune(min)) + " characters)")
		} else if l > max {
			return errors.New("value is too long (maximum is " + string(rune(max)) + " characters)")
		}
		return nil
	}
}
