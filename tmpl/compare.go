package tmpl

import "errors"
import "io"

type CompareError struct {
	ErrorMessage string
	File         string
}

func (c *CompareError) Error() string {
	return "comparasion failed for " + c.File + " due to " + c.ErrorMessage
}

func (fsRoot *fsRoot) Compare(env Environment, target Target) (Environment, error) {

	files, err := fsRoot.List(FilterFileAllowAll)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		tmpl, err := fsRoot.Reader(file)
		if err != nil {
			return nil, err
		}

		actual, err := fsRoot.Reader(file)
		if err != nil {
			if err == notFound {
				return nil, &CompareError{"target missing", file}
			} else {
				return nil, err
			}
		}

		singleEnv, err := compareSingle(env, tmpl, actual)
		if err != nil {
			return nil, err
		}

		env, err = compareJoinEnv(env, singleEnv)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func compareSingle(env Environment, ref, cmp io.Reader) (Environment, error) {
	return nil, errors.New("not implemented")
}

func compareJoinEnv(common, single Environment) (Environment, error) {
	return nil, errors.New("not implemented")
}
