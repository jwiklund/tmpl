package tmpl

import "testing"
import "reflect"
import "strings"
import "github.com/blang/vfs/memfs"

func TestEnvReadArg(t *testing.T) {
	key, value, err := envReadArg("key=value")
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

func TestEnvReadEmptyProp(t *testing.T) {
	key, value := readEnvProp("key")
	if key != "key" {
		t.Fatal(key, "not key")
	}
	if value != "" {
		t.Fatal(value, "not ''")
	}
}

func TestEnvReadPropWithdesc(t *testing.T) {
	key, value := readEnvProp("key should be key")
	if key != "key" {
		t.Fatal(key, "not key")
	}
	if value != "" {
		t.Fatal(value, "not ''")
	}
}

func TestEnvReadPropWithDefault(t *testing.T) {
	key, value := readEnvProp("key should be [default default]")
	if key != "key" {
		t.Fatal(key, "not key")
	}
	if value != "default" {
		t.Fatal(value, "not default")
	}
}

func TestEnvPropsBasic(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key")

	props, err := envReadProps(source)
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

func TestEnvNoArgs(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}

	env, err := source.GetEnvironment()
	if err != nil {
		t.Fatal(err)
	}

	if len(env.Replacements()) != 0 {
		t.Fatal("wrong replacements", env.Replacements())
	}
}

func TestEnvOneArg(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key")

	env, err := source.GetEnvironment("key=value")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env.Replacements(), []string{"key", "value"}) {
		t.Log(env.Replacements())
		t.Fatal("wrong replacements")
	}
}

func TestEnvArgWithDefault(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key [default value]\nother something else")

	env, err := source.GetEnvironment("other=value")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(env.Replacements(), []string{"key", "value", "other", "value"}) {
		t.Log(env.Replacements())
		t.Fatal("wrong replacements")
	}
}

func TestEnvMissingArg(t *testing.T) {
	source := &fsRoot{memfs.Create(), "/source"}
	testWrite(t, source, ".template", "key")

	_, err := source.GetEnvironment()
	if err == nil {
		t.Fatal("Expected missing property")
	}

	if strings.Index(err.Error(), "key") == -1 {
		t.Fatal("Wrong error, expected missing 'key'", err)
	}
}
