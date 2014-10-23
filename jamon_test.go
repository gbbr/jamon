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
BASE_URL=http://my.url
host=localhost
port=22

[routes]
users=${BASE_URL}/users
register=${BASE_URL}/register

[smtp]
connect=${host}:${port}
pipe=${connect}/pipe`,
			expected: Config{
				rootGroup: Group{
					"BASE_URL": "http://my.url",
					"host":     "localhost",
					"port":     "22",
				},
				"routes": Group{
					"users":    "http://my.url/users",
					"register": "http://my.url/register",
				},
				"smtp": Group{
					"connect": "localhost:22",
					"pipe":    "localhost:22/pipe",
				},
			},
		}, {
			contents: `
[defaults]

key=value
key2=value2`,
			expected: Config{
				"defaults": Group{
					"key":  "value",
					"key2": "value2",
				},
			},
		}, {
			contents: `
other=/sun

[defaults]
subst=my
key=${subst}${other}/val # no problem
key2=value2
key3=12${3}4`,
			expected: Config{
				rootGroup: Group{
					"other": "/sun",
				},
				"defaults": Group{
					"subst": "my",
					"key":   "my/sun/val",
					"key2":  "value2",
					"key3":  "12${3}4",
				},
			},
		}, {
			contents: `
[defaults]
subst=my
key=${subst}/val
key5=${subst}/val2
key2=value2
key3=12${3}4`,
			expected: Config{
				"defaults": Group{
					"subst": "my",
					"key":   "my/val",
					"key2":  "value2",
					"key3":  "12${3}4",
					"key5":  "my/val2",
				},
			},
		}, {
			contents: `
[defaults]
key=something
key2=this=${key} spaces`,
			expected: Config{
				"defaults": Group{
					"key":  "something",
					"key2": "this=something spaces",
				},
			},
		}, {
			contents: `
subst=my

[defaults]
key=${subst}/val
key2=value2`,
			expected: Config{
				rootGroup: Group{
					"subst": "my",
				},
				"defaults": Group{
					"key":  "my/val",
					"key2": "value2",
				},
			},
		}, {
			contents: `
ip=127.0.0.1
po.rt=23
address=${ip}:${po.rt}

[service]
address=${ip}:222
dest=${address}/get

[service2]
address=${address}/new
override=${address}/local

[service3]
reset=${address}/rset`,
			expected: Config{
				rootGroup: Group{
					"ip":      "127.0.0.1",
					"po.rt":   "23",
					"address": "127.0.0.1:23",
				},
				"service": Group{
					"address": "127.0.0.1:222",
					"dest":    "127.0.0.1:222/get",
				},
				"service2": Group{
					"address":  "127.0.0.1:23/new",
					"override": "127.0.0.1:23/new/local",
				},
				"service3": Group{
					"reset": "127.0.0.1:23/rset",
				},
			},
		}, {
			contents: `
subst=my

[defaults]
subst=priority
key=my/${subst}/val
key2=value2`,
			expected: Config{
				rootGroup: Group{
					"subst": "my",
				},
				"defaults": Group{
					"subst": "priority",
					"key":   "my/priority/val",
					"key2":  "value2",
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
				rootGroup: Group{
					"floating": "pairs",
					"of":       "keys",
				},
				"custom": Group{
					"a": "b",
					"c": "d",
				},
				"defaults": Group{
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

func TestJamon_Load_Error(t *testing.T) {
	_, err := Load("no-way-this-file-is-there.123")
	if err == nil {
		t.Error("Was expecting an error when opening an odd file")
	}
}

func TestJamon_parseLine(t *testing.T) {
	testSuite := []struct {
		input      string
		isGroup    bool
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
		isGroup, value, key, skip := parseLine(test.input)

		if isGroup != test.isGroup || value != test.value || key != test.key || skip != test.hasError {
			t.Errorf("On '%s' expected category '%t', value '%s', key '%s', but got: "+
				"'%t', '%s', '%s' and err '%s'", test.input, test.isGroup, test.value, test.key,
				isGroup, value, key, skip)
		}
	}
}

func TestJamon_Getters(t *testing.T) {
	testConfig := &Config{
		rootGroup: Group{
			"A": "B",
			"C": "D",
		},
		"category.a": Group{
			"key":  "value",
			"key2": "value2",
		},
		"category.b": Group{
			"k":   "v",
			"k.2": "v.2",
		},
	}

	// Group getters
	equality(t, testConfig.Group("category.a"), Group{
		"key":  "value",
		"key2": "value2",
	})

	equality(t, testConfig.Group("category.b"), Group{
		"k":   "v",
		"k.2": "v.2",
	})

	// Inexistent categories and keys are empty values
	equality(t, reflect.TypeOf(testConfig.Group("inexistent")).Name(), "Group")
	equality(t, len(testConfig.Group("inexistent")), 0)
	equality(t, testConfig.Get("inexistent_key"), "")
	equality(t, testConfig.Group("inexistent_cat").Get("inexistent_key"), "")
	equality(t, testConfig.Group("category.b").Get("inexistent_key"), "")

	// Default category value getters
	equality(t, testConfig.Get("A"), "B")
	equality(t, testConfig.Get("C"), "D")

	// Group value getters
	equality(t, testConfig.Group("category.a").Get("key"), "value")
	equality(t, testConfig.Group("category.a").Get("key2"), "value2")
	equality(t, testConfig.Group("category.a").Get("key"), "value")
	equality(t, testConfig.Group("category.a").Get("key2"), "value2")

	// Has functions
	equality(t, true, testConfig.Has("A"))
	equality(t, false, testConfig.Has("X"))
	equality(t, true, testConfig.HasGroup("category.a"))
	equality(t, true, testConfig.HasGroup("category.b"))
	equality(t, false, testConfig.HasGroup("category.X"))

	equality(t, true, testConfig.Group("category.b").Has("k.2"))
	equality(t, false, testConfig.Group("category.b").Has("k.X"))
	equality(t, false, testConfig.Group("category.X").Has("k.X"))
}

// Tests for deep equality
func equality(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Values %+v and %+v not equal.", a, b)
	}
}
