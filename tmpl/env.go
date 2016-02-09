package tmpl

import "bufio"
import "errors"
import "strings"

func (tmpl *fsRoot) getEnvironment(args ...string) (Environment, error) {
	return getEnvironment(false, args)
}

func (tmpl *fsRoot) getProperties(args ...string) (Environment, error) {
	return getEnvironment(true, args)
}

func (tmpl *fsRoot) getEnvironment(ignoreMissing bool, args ...string) (Environment, error) {
	props, err := envReadProps(tmpl)
	if err != nil {
		return nil, err
	}
	env := map[string]string{}

	for _, arg := range args {
		key, value, err := envReadArg(arg)
		if err != nil {
			return nil, err
		}
		_, ok := props[key]
		if !ok {
			return nil, errors.New("undefined property " + key)
		}
		env[key] = value
	}

	for key, value := range props {
		_, ok := env[key]
		if ok {
			continue
		}
		if len(value) > 0 {
			env[key] = value
		} else if !ignoreMissing {
			return nil, errors.New("missing required property " + key)
		}
	}
	return Environment(env), nil
}

func envReadProps(tmpl *fsRoot) (map[string]string, error) {
	reader, err := tmpl.Reader(".template")
	if err == notFound {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	props := map[string]string{}
	s := bufio.NewScanner(reader)
	for s.Scan() {
		key, value := readEnvProp(s.Text())
		props[key] = value
	}
	if err = s.Err(); err != nil {
		return nil, err
	}
	return props, nil
}

func envReadArg(arg string) (string, string, error) {
	eq := strings.Index(arg, "=")
	if eq == -1 {
		return "", "", errors.New("invalid argument '" + arg + "'")
	}
	start := 0
	for start < len(arg) && arg[start] == '-' {
		start = start + 1
	}
	return arg[start:eq], arg[eq+1 : len(arg)], nil
}

func readEnvProp(prop string) (string, string) {
	space := strings.Index(prop, " ")
	if space == -1 {
		return prop, ""
	}
	def := strings.Index(prop[space:], "[default ")
	if def == -1 {
		return prop[0:space], ""
	}
	def = space + def + len("[default ")
	enddef := strings.Index(prop[def:], "]")
	if enddef == -1 {
		enddef = len(prop)
	} else {
		enddef = def + enddef
	}
	return prop[0:space], prop[def:enddef]
}
