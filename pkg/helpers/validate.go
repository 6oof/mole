package helpers

import (
	"regexp"
)

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*(\.[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*)*\.[a-zA-Z]{2,}$`

const domainRegex = `^(\*\.)?([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z\p{L}]{2,}(:\d+)?$`

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateCaddyDomain(domain string) bool {
	re := regexp.MustCompile(domainRegex)
	return re.MatchString(domain)
}
