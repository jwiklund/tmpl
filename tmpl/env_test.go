package tmpl

import "testing"
import "reflect"
import "strings"
import "github.com/blang/vfs/memfs"

func TestNoArgs(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}

	env, err := GetEnvironment(source)
	if err != nil {
		t.Fatal(err)
	}

	if len(env.Replacements()) != 0 {
		t.Fatal("wrong replacements", env.Replacements())
	}
}

func TestBasicProps(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key")

	props, err := readEnvProps(source)
	if err != nil {
		t.Fatal(err)
	}

	if len(props) != 1 {
		t.Fatal("expected one entry", props)
	}

	if props["key"] != "" {
		t.Fatal("expected key=\"\"", props["key"])
	}
}

func TestReadEnvArg(t *testing.T) {
	key, value, err := readEnvArg("key=value")
	if err != nil {
		t.Fatal(err)
	}
	if key != "key" {
		t.Fatal(err)
	}
	if value != value {
		t.Fatal(value)
	}
}

func TestReadEmptyProp(t *testing.T) {
	key, value := readEnvProp("key")
	if key != "key" {
		t.Fatal(key, "not key")
	}
	if value != "" {
		t.Fatal(value, "not ''")
	}
}

func TestReadPropWithdesc(t *testing.T) {
	key, value := readEnvProp("key should be key")
	if key != "key" {
		t.Fatal(key, "not key")
	}
	if value != "" {
		t.Fatal(value, "not ''")
	}
}

func TestReadPropWithDefault(t *testing.T) {
	key, value := readEnvProp("key should be [default default]")
	if key != "key" {
		t.Fatal(key, "not key")
	}
	if value != "default" {
		t.Fatal(value, "not default")
	}
}

func TestOneArg(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key")

	env, err := GetEnvironment(source, "key=value")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env.Replacements(), []string{"key", "value"}) {
		t.Log(env.Replacements())
		t.Fatal("wrong replacements")
	}
}

func TestArgWithDefault(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key [default value]\nother something else")

	env, err := GetEnvironment(source, "other=value")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env.Replacements(), []string{"key", "value", "other", "value"}) {
		t.Log(env.Replacements())
		t.Fatal("wrong replacements")

	}
}

func TestMissingArg(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key")

	_, err := GetEnvironment(source)
	if err == nil {
		t.Fatal("Expected missing property")
	}

	if strings.Index(err.Error(), "key") == -1 {
		t.Fatal("Wrong error, expected missing 'key'", err)
	}
}
