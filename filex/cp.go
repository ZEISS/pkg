package filex

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// CopyFile ...
func CopyFile(src, dst string, mkdir bool) (int64, error) {
	src, err := AbsolutePath(src)
	if err != nil {
		return 0, err
	}

	dst, err = AbsolutePath(dst)
	if err != nil {
		return 0, err
	}

	if mkdir {
		err = MkdirAll(filepath.Dir(dst), 0o755)
		if err != nil {
			return 0, err
		}
	}

	return copy2(src, dst)
}

// CopyEmbedFile ...
func CopyEmbedFile(src, dst string, fs embed.FS) (int64, error) {
	var size int64

	f, err := fs.ReadFile(src)
	if err != nil {
		return size, err
	}

	dest, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return size, err
	}
	defer dest.Close()

	size, err = io.Copy(dest, bytes.NewReader(f))
	if err != nil {
		return size, err
	}

	return size, nil
}

// AbsolutePath ...
func AbsolutePath(path string) (string, error) {
	path, err := ExpandHomeFolder(path)
	if err != nil {
		return "", err
	}

	return filepath.Abs(path)
}

// ExpandHomeFolder ...
func ExpandHomeFolder(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	var buffer bytes.Buffer
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	_, err = buffer.WriteString(usr.HomeDir)
	if err != nil {
		return "", err
	}

	_, err = buffer.WriteString(strings.TrimPrefix(path, "~"))
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// PrependHomeFolder ...
func PrependHomeFolder(path string) (string, error) {
	if strings.HasPrefix(path, string(os.PathSeparator)) {
		return path, nil
	}

	var buffer bytes.Buffer
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	_, err = buffer.WriteString(filepath.Join(usr.HomeDir, filepath.Clean(path)))
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// Transformer ...
type Transformer func(path string) (string, error)

// PathTransform ...
func PathTransform(path string, funcs ...Transformer) (string, error) {
	p := path

	var err error
	for _, fn := range funcs {
		p, err = fn(path)
		if err != nil {
			return "", err
		}
	}

	return p, nil
}

func copy2(src, dst string) (int64, error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sfi.Mode().IsRegular() {
		return 0, fmt.Errorf("copyfile: %s is not a regular file", src)
	}

	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return 0, err
	}
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() { _ = destination.Close() }()

	n, err := io.Copy(destination, source)

	return n, err
}
