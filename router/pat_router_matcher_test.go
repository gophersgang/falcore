package router

import (
	"reflect"
	"testing"
)

var matcherTestData = []struct {
	name     string
	pattern  string
	path     string
	match    bool
	captures map[string]string
}{
	{
		"slash",
		"/",
		"/",
		true,
		map[string]string{},
	},
	{
		"slash mismatch",
		"/",
		"/foo",
		false,
		nil,
	},
	{
		"basic",
		"/foo",
		"/foo",
		true,
		map[string]string{},
	},
	{
		"basic mismatch",
		"/foo",
		"/bar",
		false,
		nil,
	},
	{
		"basic underrun",
		"/foo",
		"/",
		false,
		nil,
	},
	{
		"basic overrun",
		"/foo",
		"/foo/bar",
		false,
		nil,
	},
	{
		"capture",
		"/:foo",
		"/bar",
		true,
		map[string]string{"foo": "bar"},
	},
}

func Test_pattern_match(t *testing.T) {
	for _, test := range matcherTestData {
		pat, _ := ParsePattern(test.pattern)
		match, captures := pat.match(test.path)
		if test.match != match {
			t.Errorf("[%v] Match result mismatch.  Expected %v. Got %v.", test.name, test.match, match)
		}
		if !reflect.DeepEqual(captures, test.captures) {
			t.Errorf("[%v] Captures don't match. Expected %v. Got %v", test.name, test.captures, captures)
		}
	}

}
