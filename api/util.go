package main

import "regexp"

func isEmailAddressRegex(email string) bool {
	// Define a regex pattern for validating an email address
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Match the email against the regex pattern
	return emailRegex.MatchString(email)
}
