package tmpl

import "testing"
import "os"
import "github.com/blang/vfs/memfs"
import "bytes"
import "reflect"

func TestIOWriteFile(t *testing.T) {
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

func TestIOReadFile(t *testing.T) {
	fs := memfs.Create()
	fs.Mkdir("/root", 0775)
	file, err := fs.OpenFile("/root/file", os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		t.Fatal("could not create test file", err)
	}

	_, err = file.Write([]byte("hello"))
	if err != nil {
		t.Fatal("could not write test file", err)
	}

	target := fsRoot{fs, "/root"}
	reader, err := target.Reader("file")
	if err != nil {
		t.Fatal("could not create reader")
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		t.Fatal("failed reading", err)
	}

	if buf.String() != "hello" {
		t.Log(buf.String())
		t.Fatal("expected hello")
	}
}

func TestIOReadMissing(t *testing.T) {
	target := fsRoot{memfs.Create(), "/root"}
	_, err := target.Reader("file")

	if err != notFound {
		t.Fatal("expected not found", err)
	}
}

func TestIOListEmpty(t *testing.T) {
	fs := memfs.Create()
	fs.Mkdir("/root", 0775)
	target := fsRoot{fs, "/root"}
	list, err := target.List(FilterFileAllowAll)

	if err != nil {
		t.Fatal("list failed", err)
	}

	if len(list) != 0 {
		t.Fatal("expected no entries", list)
	}
}

func TestIOList(t *testing.T) {
	target := &fsRoot{memfs.Create(), "/root"}
	testWrite(t, target, "sub1", "hello1")
	testWrite(t, target, "sub2", "hello2")

	list, err := target.List(FilterFileAllowAll)

	if err != nil {
		t.Fatal("list failed", err)
	}

	if !reflect.DeepEqual(list, []string{"/sub1", "/sub2"}) {
		t.Fatal("wrong list", list)
	}
}

func TestIOListSubs(t *testing.T) {
	target := &fsRoot{memfs.Create(), "/root"}
	testWrite(t, target, "sub1/leaf1", "hello1")
	testWrite(t, target, "sub1/leaf2", "hello2")
	testWrite(t, target, "sub2/leaf3", "hello3")

	list, err := target.List(FilterFileAllowAll)

	if err != nil {
		t.Fatal("list failed", err)
	}

	if !reflect.DeepEqual(list, []string{"/sub1/leaf1", "/sub1/leaf2", "/sub2/leaf3"}) {
		t.Fatal("wrong list", list)
	}
}

func TestIOListRead(t *testing.T) {
	target := &fsRoot{memfs.Create(), "/root"}
	testWrite(t, target, "sub1", "hello")

	list, err := target.List(FilterFileAllowAll)

	if err != nil {
		t.Fatal("list failed", err)
	}

	if len(list) != 1 {
		t.Fatal("wrong list", list)
	}

	reader, err := target.Reader(list[0])
	if err != nil {
		t.Fatal(err)
	}

	var bytes bytes.Buffer
	bytes.ReadFrom(reader)
	if bytes.String() != "hello" {
		t.Log(bytes.String())
		t.Fatal("wrong result")
	}
}

func TestIORestrict(t *testing.T) {
	target := &fsRoot{memfs.Create(), "/"}
	testWrite(t, target, "sub1/leaf1", "hello")

	restricted := target.Restrict("sub1")
	if testRead(t, restricted, "leaf1") != "hello" {
		t.Fatal("expected hello")
	}
}
