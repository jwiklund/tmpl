package tmpl

import "bytes"
import "testing"

func testWrite(t *testing.T, root *fsRoot, name, value string) {
	target, err := root.Writer(name)
	if err != nil {
		t.Fatal("could not create writer", err)
	}
	if _, err := target.Write([]byte(value)); err != nil {
		t.Fatal("could not write", value)
	}
}

func testRead(t *testing.T, root *fsRoot, name string) string {
	source, err := root.Reader(name)
	if err != nil {
		t.Fatal("could not create reader", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(source)
	if err != nil {
		t.Fatal("Could not read", err)
	}
	t.Log("read", name, "contents", buf.String())
	return buf.String()
}
