package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"asdf", false},
		{"aa@aa.com", true},
		{"user.name+tag@sub.domain.com", true},
		{"@missingusername.com", false},
		{"username@.com", false},
		{"username@domain..com", false},
		{"username@domain.com.", false},
		{"username@domain.c", false},
		{"user@domain_with_space .com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			assert.Equal(t, tt.expected, ValidateEmail(tt.email))
		})
	}
}

func TestValidateCaddyDomain(t *testing.T) {
	tests := []struct {
		domain   string
		expected bool
	}{
		{"example.com", true},
		{"*.example.com", true},
		{"subdomain.example.com", true},
		{"example", false},
		{"example..com", false},
		{"example.com:8080", true},
		{"*.example..com", false},
		{"example@domain.com", false},
		{"example.com:port", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			assert.Equal(t, tt.expected, ValidateCaddyDomain(tt.domain))
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"project1", true},
		{"project-name", true},
		{"project_name", true},
		{"project123", true},
		{"123project", true},
		{"project-123", true},
		{"project_123", true},
		{"project@", false},
		{"project name", false},
		{"ProjectName", false},
		{"project.name", false},
		{"project-name-", false},
		{"-project", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ValidateProjectName(tt.name))
		})
	}
}
