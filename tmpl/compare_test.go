package tmpl

import "testing"
import "github.com/blang/vfs/memfs"

func TestCompareSelfNoProperties(t *testing.T) {
	fs := &fsRoot{memfs.Create(), "/"}
	testWrite(t, fs, "test", "value")

	env, err := fs.Compare(nil, fs)
	if err != nil {
		t.Fatal(err)
	}

	if len(env) != 0 {
		t.Fatal("Expected no env", env)
	}
}
