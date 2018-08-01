// +build gofuzz

package validation

import "github.com/go-ozzo/ozzo-validation/is"

func Fuzz(data []byte) int {
	rules := []Rule{
		Required,
		Length(5, 100),
		is.Email,
		is.DNSName,
		is.Port,
		is.IPv4,
	}

	if err := Validate(string(data), rules...); err != nil {
		return 0
	}
	return 1
}
