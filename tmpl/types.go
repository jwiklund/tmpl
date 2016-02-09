package tmpl

import "io"
import "errors"
import "github.com/blang/vfs"

var notFound error = errors.New("does not exist")

type Target interface {
	Writer(path string) (io.Writer, error)
	Reader(path string) (io.Reader, error)
}

type Template interface {
	Create(target Target, args ...string) error

	Compare(target Target, args ...string) (Environment, error)
}

type Environment map[string]string

func NewTemplate(path string) (Template, error) {
	return &fsRoot{vfs.OS(), path}, nil
}

func GetTarget(path string) (Target, error) {
	return &fsRoot{vfs.OS(), path}, nil
}

type fsRoot struct {
	fs   vfs.Filesystem
	root string
}
