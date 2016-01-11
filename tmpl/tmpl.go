package tmpl

import "errors"
import "golang.org/x/tools/godoc/vfs"

type Template interface {
	Create() error
}

func New(path string) (Template, error) {
	return &fsTemplate{vfs.OS(path)}, nil
}

type fsTemplate struct {
	fs vfs.FileSystem
}

func (*fsTemplate) Create() error {
	return errors.New("not implemented")
}
