package jamon

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestJamon_Load(t *testing.T) {
	testSuite := []struct {
		contents string
		expected Config
	}{
		{
			contents: `
[defaults]

key=value
key2=value2`,
			expected: Config{
				"defaults": Category{
					"key":  "value",
					"key2": "value2",
				},
			},
		}, {
			contents: `
floating=pairs
of=keys # with a comment

# These are the defaults
[defaults]

key=value

# This key is important
key2=value2

[custom] # This is the best category
a=b
c=d`,
			expected: Config{
				"JAMON.NO_CATEGORY": Category{
					"floating": "pairs",
					"of":       "keys",
				},
				"custom": Category{
					"a": "b",
					"c": "d",
				},
				"defaults": Category{
					"key":  "value",
					"key2": "value2",
				},
			},
		},
	}

	for _, test := range testSuite {
		err := ioutil.WriteFile("testfile.tmp", []byte(test.contents), 0777)
		if err != nil {
			t.Errorf("Error creating temporary file: %s", err)
		}

		config, err := Load("testfile.tmp")
		if err != nil {
			t.Error("Error loading file: %s", err)
		}

		if !reflect.DeepEqual(config, test.expected) {
			t.Errorf("Expected %+v, got %+v", test.expected, config)
		}

		err = os.Remove("testfile.tmp")
		if err != nil {
			t.Errorf("Error removing file: %s", err)
		}
	}
}

func TestJamon_parseLine(t *testing.T) {
	testSuite := []struct {
		input      string
		isCategory bool
		value, key string
		hasError   bool
	}{
		// Categories
		{"[category]", true, "category", "", false},
		{"[category.name] # weird comment", true, "category.name", "", false},

		// Errors
		{"category", false, "", "", true},

		// Key / value pairs
		{"key=value", false, "value", "key", false},
		{"\tkey=value=key", false, "value=key", "key", false},
		{" key=value   # with comment", false, "value", "key", false},
		{"key=value # with comment # inception", false, "value", "key", false},

		// Empty lines and dodgy chars
		{"\r\n", false, "", "", true},
		{"\n", false, "", "", true},
		{"", false, "", "", true},
		{"     ", false, "", "", true},
		{"\t     ", false, "", "", true},

		// Comments
		{"# Line with comment", false, "", "", true},
		{"## Line with comment # asd", false, "", "", true},
	}

	for _, test := range testSuite {
		isCategory, value, key, skip := parseLine(test.input)

		if isCategory != test.isCategory || value != test.value || key != test.key || skip != test.hasError {
			t.Errorf("On '%s' expected category '%t', value '%s', key '%s', but got: "+
				"'%t', '%s', '%s' and err '%s'", test.input, test.isCategory, test.value, test.key,
				isCategory, value, key, skip)
		}
	}
}
