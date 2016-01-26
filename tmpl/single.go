package tmpl

import "io"
import "strings"
import "bufio"
import "errors"

func (e Environment) Replacements() []string {
	replacements := make([]string, 2*len(e))
	i := 0
	for key, value := range e {
		replacements[i] = key
		replacements[i+1] = value
		i = i + 2
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
