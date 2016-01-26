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

func TestSingleProp(t *testing.T) {
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
	err := WriteSingle(Environment(props), r, &o)
	if err != nil {
		t.Fatal("Could not replace single ", err.Error())
	}
	return o.String()
}

func TestSingleNoReplacement(t *testing.T) {
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

func TestSingleSubStringReplacement(t *testing.T) {
	res := replace(t, "hello PROPERTY1", "PROPERTY1=property1")
	if res != "hello property1\n" {
		t.Fatal("Expected hello property1")
	}
}

func TestSingleMultipleReplacements(t *testing.T) {
	res := replace(t, "hello PROP1 and PROP2", "PROP1=t1,PROP2=t2")
	if res != "hello t1 and t2\n" {
		t.Fatal("Expecting hello t1 and t2 but got", res)
	}
}
