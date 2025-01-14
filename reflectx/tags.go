package reflectx

import (
	"strings"
	"unicode"

	"github.com/zeiss/pkg/utilx"
)

// TagOptions ...
type TagOptions string

// ParseTag splits a struct field's json tag into its name and
// comma-separated options.
func ParseTag(tag string) (string, TagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, TagOptions(opt)
}

// IsValidTag returns true if the tag is not empty.
func IsValidTag(tag string) bool {
	if utilx.Empty(tag) {
		return false
	}

	for _, c := range tag {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}

	return true
}
