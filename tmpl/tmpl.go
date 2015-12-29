package tmpl

import "errors"

type Template interface {
	Create() error
}

func New(path string) (Template, error) {
	return &pathTemplate{path}, nil
}

type pathTemplate struct {
	path string
}

func (*pathTemplate) Create() error {
	return errors.New("not implemented")
}
