package tmpl

import "bytes"
import "testing"
import "github.com/blang/vfs/memfs"

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
	t.Log(buf.String())
	return buf.String()
}

func TestSimpleCreate(t *testing.T) {
	mem := memfs.Create()
	source := &fsRoot{mem, "/source"}
	target := &fsRoot{mem, "/target"}

	testWrite(t, source, "template", "hello")

	err := source.Create(target)
	if err != nil {
		t.Fatal("Template creation failed", err)
	}

	if testRead(t, target, "template") != "hello\n" {
		t.Fatal("Expected hello")
	}
}
