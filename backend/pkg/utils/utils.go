package utils

import "regexp"

const emailRegexPattern = `^[^\s@]+@[^\s@]+\.[^\s@]+$`

var emailRegex = regexp.MustCompile(emailRegexPattern)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
