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
	GetEnvironment(args ...string) (Environment, error)

	Create(env Environment, target Target) error
	Compare(env Environment, target Target) (Environment, error)
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
