package tmpl

import "io"
import "errors"
import "github.com/blang/vfs"

var NOT_FOUND error = errors.New("does not exist")

type Target interface {
	Writer(path string) (io.Writer, error)
}

type Template interface {
	Properties() (io.Reader, error)
	Create(env Environment, target Target) error
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
