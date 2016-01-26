package tmpl

import "testing"
import "os"
import "github.com/blang/vfs/memfs"
import "bytes"

func TestWriteFile(t *testing.T) {
	fs := memfs.Create()
	target := fsRoot{fs, "/root"}

	writer, err := target.Writer("file")
	if err != nil {
		t.Fatal("Could not get writer", err)
	}

	_, err = writer.Write([]byte("hello"))
	if err != nil {
		t.Fatal("Could not write", err)
	}

	file, err := fs.OpenFile("/root/file", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal("File not created", err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		t.Fatal("Failed reading", err)
	}

	if buf.String() != "hello" {
		t.Log(buf.String())
		t.Fatal("expected hello")
	}
}

func TestReadFile(t *testing.T) {
	fs := memfs.Create()
	fs.Mkdir("/root", 0775)
	file, err := fs.OpenFile("/root/file", os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		t.Fatal("Could not create test file", err)
	}

	_, err = file.Write([]byte("hello"))
	if err != nil {
		t.Fatal("Could not write test file", err)
	}

	target := fsRoot{fs, "/root"}
	reader, err := target.Reader("file")
	if err != nil {
		t.Fatal("Could not create reader")
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		t.Fatal("Failed reading", err)
	}

	if buf.String() != "hello" {
		t.Log(buf.String())
		t.Fatal("Expected hello")
	}
}

func TestReadMissing(t *testing.T) {
	target := fsRoot{memfs.Create(), "/root"}
	_, err := target.Reader("file")

	if err != NOT_FOUND {
		t.Fatal("Expected not found", err)
	}
}
