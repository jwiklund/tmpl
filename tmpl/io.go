package tmpl

import "errors"
import "io"
import "os"
import "path/filepath"
import "sort"
import "github.com/blang/vfs"

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

func (fsRoot *fsRoot) Properties() (io.Reader, error) {
	return fsRoot.Reader(".template")
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

type FilterFile func(string, bool) bool

func FilterFileAllowAll(path string, dir bool) bool {
	return true
}

func (fsRoot *fsRoot) List(filter FilterFile) ([]string, error) {
	root := fsRoot.root
	if root[len(root)-1] != '/' {
		root = root + string(fsRoot.fs.PathSeparator())
	}

	return listCollect(fsRoot.fs, root, "/", filter)
}

type byName []os.FileInfo

func (l byName) Len() int {
	return len(l)
}

func (l byName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l byName) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}

func listCollect(fs vfs.Filesystem, realPath, virtualPath string, filter FilterFile) ([]string, error) {
	files := []string{}
	contents, err := fs.ReadDir(realPath)
	if err != nil {
		return nil, err
	}
	sort.Sort(byName(contents))
	for _, content := range contents {
		if content.IsDir() {
			if filter(virtualPath+content.Name(), true) {
				subfiles, err := listCollect(fs, realPath+content.Name()+"/", virtualPath+content.Name()+"/", filter)
				if err != nil {
					return nil, err
				}
				files = append(files, subfiles...)
			}
		} else {
			if filter(virtualPath+content.Name(), false) {
				files = append(files, virtualPath+content.Name())
			}
		}
	}
	return files, nil
}
