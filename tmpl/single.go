package tmpl

import "io"
import "strings"
import "bufio"
import "errors"

func WriteSingle(in io.Reader, out io.Writer, properties map[string]string) error {
	replacements := make([]string, 2*len(properties))
	i := 0
	for key, value := range properties {
		replacements[i] = key
		replacements[i+1] = value
		i = i + 1
	}
	replacer := strings.NewReplacer(replacements...)

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
