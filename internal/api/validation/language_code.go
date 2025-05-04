package validation

import (
	"regexp"
)

var re = regexp.MustCompile(`^[a-zA-Z]{2}(-[a-zA-Z]{2})?$`)

func IsValidLanguageCode(code string) bool {
	return re.MatchString(code)
}
