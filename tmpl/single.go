package tmpl

import "io"
import "sort"
import "strings"
import "bufio"
import "errors"

func (e Environment) Replacements() []string {
	keys := []string{}
	for key, _ := range e {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	replacements := make([]string, 2*len(keys))
	for i, key := range keys {
		replacements[i*2] = key
		replacements[i*2+1] = e[key]
	}
	return replacements
}

func WriteSingle(env Environment, in io.Reader, out io.Writer) error {
	replacer := strings.NewReplacer(env.Replacements()...)

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		replacer.WriteString(out, scanner.Text())
		out.Write([]byte("\n"))
	}

	if err := scanner.Err(); err != nil {
		return errors.New("Failed to read template: " + err.Error())
	}
	return nil
}
