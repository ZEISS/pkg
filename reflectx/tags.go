package reflectx

import "strings"

// TagOptions ...
type TagOptions string

// ParseTag splits a struct field's json tag into its name and
// comma-separated options.
func ParseTag(tag string) (string, TagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, TagOptions(opt)
}
