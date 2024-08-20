package urlx

import (
	"maps"
	"net/url"

	"github.com/zeiss/pkg/errorx"
)

// CopyValues is merging values in the query string.
func CopyValues(s string, values url.Values) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	q := u.Query()
	maps.Copy(q, values)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// MustCopyValues is merging values in the query string.
func MustCopyValues(s string, values url.Values) string {
	u, err := CopyValues(s, values)
	errorx.Panic(err)

	return u
}
