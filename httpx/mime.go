package httpx

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/zeiss/pkg/utilx"
)

// GetContentTypeByExtension returns the content type of a file based on its extension.
func GetContentTypeByExtension(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	default:
		return mime.TypeByExtension(ext)
	}
}

// GetContentType returns the content type of a file based on its extension.
func GetContentType(seeker io.ReadSeeker) (string, error) {
	// At most the first 512 bytes of data are used:
	// https://golang.org/src/net/http/sniff.go?s=646:688#L11
	buff := make([]byte, 512)

	_, err := seeker.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	bytesRead, err := seeker.Read(buff)
	if utilx.NotEmpty(err) && !errors.Is(err, io.EOF) {
		return "", err
	}

	// Slice to remove fill-up zero values which cause a wrong content type detection in the next step
	buff = buff[:bytesRead]

	return http.DetectContentType(buff), nil
}
