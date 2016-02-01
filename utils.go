package vcardgen

import (
	"strings"

	"github.com/dchest/validator"
)

// Makes encoding substitutions for vcard. Maybe.
func encodeString(value string) string {
	value = strings.Replace(value, "\n", "\\n", -1)
	value = strings.Replace(value, ",", "\\,", -1)
	value = strings.Replace(value, ";", "\\;", -1)
	return value
}

// Normalises email addresses; lowercase, no flanking space, and a few other tricks.
func normaliseEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return validator.NormalizeEmail(email)
}
