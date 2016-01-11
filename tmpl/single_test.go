package tmpl

import "testing"
import "strings"
import "bytes"
import "bufio"

func toProps(properties string) map[string]string {
	p := make(map[string]string)
	propreader := bufio.NewReader(strings.NewReader(properties))
	for line, err := propreader.ReadString('='); err == nil; {
		rest, _ := propreader.ReadString(',')
		p[strings.TrimRight(line, "=")] = strings.TrimRight(rest, ",")
		line, err = propreader.ReadString('=')
	}
	return p
}

func TestProp(t *testing.T) {
	p := toProps("test=ing,test2=ing2")
	if p["test"] != "ing" {
		t.Fatal("test", p)
	}
	if p["test2"] != "ing2" {
		t.Fatal("test2", p)
	}
}

func replace(t *testing.T, input, properties string) string {
	r := strings.NewReader(input)
	var o bytes.Buffer
	props := toProps(properties)
	err := WriteSingle(r, &o, props)
	if err != nil {
		t.Fatal("Could not replace single ", err.Error())
	}
	return o.String()
}

func TestNoReplacement(t *testing.T) {
	res := replace(t, "test", "")
	if res != "test\n" {
		t.Fatal("Expected 'test' but got '" + res + "'")
	}
}

func TestSingleReplacement(t *testing.T) {
	res := replace(t, "test", "test=ing")
	if res != "ing\n" {
		t.Fatal("Expected ing but got '" + res + "'")
	}
}

func TestSingleMultilineReplacement(t *testing.T) {
	res := replace(t, "test\ntest2", "test=ing")
	if res != "ing\ning2\n" {
		t.Fatal("Expected ing\ning2 but got '" + res + "'")
	}
}

func TestSingleMultilineReplacementEOL(t *testing.T) {
	res := replace(t, "test\ntest2\n", "test=ing")
	if res != "ing\ning2\n" {
		t.Fatal("Expected ing\ning2 but got '" + res + "'")
	}
}
