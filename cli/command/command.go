package command

import "strings"

// validateIntGreaterThan returns an integer validator.
//
// This validator ensures the value received is greater than the given minimum.
func validateIntGreaterThan(min int) func(int) (int, bool) {
	return func(val int) (int, bool) {
		return val, val > min
	}
}

// validateStringIsOneOf returns a string validator.
//
// This validator ensures the value received is in the given list of valid values.
func validateStringIsOneOf(valid []string) func(string) (string, bool) {
	return func(val string) (string, bool) {
		for _, v := range valid {
			if val == v {
				return val, true
			}
		}
		return val, false
	}
}

var transformStringToCaps = func(val string) (string, bool) {
	return strings.ToUpper(val), true
}
