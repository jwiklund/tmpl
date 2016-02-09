package tmpl

import "errors"
import "io"
import "bufio"

type CompareError struct {
	ErrorMessage string
	File         string
}

func (c *CompareError) Error() string {
	return "comparasion failed for " + c.File + " due to " + c.ErrorMessage
}

func (fsRoot *fsRoot) Compare(target Target, args ...string) (Environment, error) {

	refEnv, err := fsRoot.getProperties(args)
	if err != nil {
		return nil, err
	}
	env := refEnv

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

		singleEnv, err := compareSingle(refEnv, tmpl, actual)
		if err != nil {
			return nil, &CompareError{err.Error(), file}
		}

		env, err = compareJoinEnv(env, singleEnv)
		if err != nil {
			return nil, &CompareError{err.Error(), file}
		}
	}

	return env, nil
}

func compareSingle(refEnv Environment, ref, cmp io.Reader) (Environment, error) {
	refLines := bufio.NewScanner(ref)
	cmpLines := bufio.NewScanner(cmp)

	var env map[string]string

	for {
		scannedRef := refLines.Scan()
		scannedCmp := cmpLines.Scan()
		if scannedRef && !scannedCmp {
			return nil, errors.New("more lines in reference")
		}
		if !scannedRef && scannedCmp {
			return nil, errors.New("more lines in comparasion")
		}
		if !scannedRef && !scannedCmp {
			return env, nil
		}

		lineEnv, err := compareLine(refEnv, env, scannedRef.Text(), scannedCmp.Text())

		if err != nil {
			return nil, err
		}

		env = compareJoinEnv(env, lineEnv)
		for key, value := range lineEnv {
			earlier, ok := env[key]
			if ok && earlier != value {
				return nil, errors.New("invalid file, both " + earlier + " and " + value + " for " + key)
			}
			env[key] = value
		}
	}
	return env, nil
}

func compareLine(refEnv, env Environment, ref, cmp string) (Environment, error) {
	return nil, errors.New("not implemented")
}

func compareJoinEnv(common, single Environment) (Environment, error) {
	return nil, errors.New("not implemented")
}
