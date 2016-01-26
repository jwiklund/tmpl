package tmpl

import "os"
import "io"
import "path/filepath"
import "errors"
import "github.com/blang/vfs"

var NOT_FOUND error = errors.New("does not exist")

type Target interface {
	Writer(path string) (io.Writer, error)
}

type Template interface {
	Create(target Target) error
}

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

func (fsRoot *fsRoot) Writer(path string) (io.Writer, error) {
	realPath := fsRoot.root + string(fsRoot.fs.PathSeparator()) + path
	dir := filepath.Dir(realPath)
	err := vfs.MkdirAll(fsRoot.fs, dir, 0775)
	if err != nil {
		return nil, errors.New("could not create dir " + dir + " due to " + err.Error())
	}
	file, err := fsRoot.fs.OpenFile(realPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return nil, errors.New("could not open file " + realPath + " due to " + err.Error())
	}
	return file, nil
}

func (fsRoot *fsRoot) Reader(path string) (io.Reader, error) {
	realPath := fsRoot.root + string(fsRoot.fs.PathSeparator()) + path
	file, err := fsRoot.fs.OpenFile(realPath, os.O_RDONLY, 0)
	if err != nil {
		pathError, isPathError := err.(*os.PathError)
		if isPathError {
			if pathError.Err.Error() == "file does not exist" {
				return nil, NOT_FOUND
			}
			return nil, errors.New("could not open file due to " + err.Error())
		}
		return nil, errors.New("could not open file " + realPath + " due to " + err.Error())
	}
	return file, nil
}
