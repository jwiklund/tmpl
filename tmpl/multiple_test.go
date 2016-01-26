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
	t.Log("read", name, "contents", buf.String())
	return buf.String()
}

func TestSimpleCreate(t *testing.T) {
	mem := memfs.Create()
	source := &fsRoot{mem, "/source"}
	target := &fsRoot{mem, "/target"}

	testWrite(t, source, "template", "hello")

	err := source.Create(Environment(nil), target)
	if err != nil {
		t.Fatal("Template creation failed", err)
	}

	if testRead(t, target, "template") != "hello\n" {
		t.Fatal("Expected hello")
	}
}

func TestProperties(t *testing.T) {
	mem := memfs.Create()
	source := &fsRoot{mem, "/source"}
	target := &fsRoot{mem, "/source"}

	testWrite(t, source, "template1", "hello PROPERTY1")
	testWrite(t, source, "template2", "hello PROPERTY2")

	err := source.Create(Environment(map[string]string{
		"PROPERTY1": "property1",
		"PROPERTY2": "property2",
	}), target)

	if err != nil {
		t.Fatal(err)
	}

	if testRead(t, target, "template1") != "hello property1\n" {
		t.Fatal("wrong property1")
	}

	if testRead(t, target, "template2") != "hello property2\n" {
		t.Fatal("wrong property2")
	}
}
