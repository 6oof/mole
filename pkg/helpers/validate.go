package helpers

import (
	"regexp"
)

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateCaddyDomain(domain string) bool {
	re := regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z\p{L}]{2,}(:\d+)?$`)
	return re.MatchString(domain)
}
