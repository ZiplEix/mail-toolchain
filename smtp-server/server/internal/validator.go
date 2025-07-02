package internal

import "regexp"

var emailRegex = regexp.MustCompile(`^<[^<>@]+@[^<>@]+\.[^<>@]+>$`)

func IsValidEmail(addr string) bool {
	return emailRegex.MatchString(addr)
}
