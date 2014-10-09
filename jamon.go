/*
Package Jamon is an INI-like configuration file parser. An example configuration
file may look like this:
	address=127.0.0.1:1234 # root-level values

	[defaults]
	key=value
	name=Gabriel

	[category]
	key=value
Trailing comments are also allowed, and root-level keys are only accepted at the
top of the file
*/
package jamon

import (
	"bufio"
	"os"
	"strings"
)

// A configuration type may hold multiple categories of settings
type Config map[string]Category

// A category holds key-value pairs of settings
type Category map[string]string

// Internal name for category that holds settings at root-level
const defaultCategory = "JAMON.NO_CATEGORY"

// Returns the value of a root-level key
func (c Config) Get(key string) string { return c[defaultCategory].Get(key) }

// Verifies the existence of a root-level key
func (c Config) HasKey(key string) bool {
	_, ok := c[defaultCategory][key]
	return ok
}

// Returns a category by name. If the category does not exist, an empty category
// is returned. This is to avoid multiple return values in order to facilitate
// chaining.
func (c Config) Category(name string) Category { return c[name] }

// Verifies if a category exists
func (c Config) HasCategory(category string) bool {
	_, ok := c[category]
	return ok
}

// Returns a key from a category
func (c Category) Get(key string) string { return c[key] }

// Verifies if the category has a key
func (c Category) HasKey(key string) bool {
	_, ok := c[key]
	return ok
}

// Loads a configuration file
func Load(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	config := Config{}
	currentCategory := defaultCategory

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		isCategory, value, key, skip := parseLine(string(line))

		switch {
		case skip:
			continue

		case isCategory:
			currentCategory = value
			continue

		case config[currentCategory] == nil:
			config[currentCategory] = make(Category)
			fallthrough
		default:
			config[currentCategory][key] = value
		}
	}

	return config, nil
}

// Attempts to parse an entry in the config file. The first return value specifies
// whether 'value' is the name of a category or the value of a key. Skip indicates
// whether the line was a comment or could not be parsed.
func parseLine(line string) (isCategory bool, value, key string, skip bool) {
	line = strings.SplitN(line, "#", 2)[0]
	line = strings.Trim(line, " \t\r")

	// Is comment or empty line?
	if len(line) == 0 {
		skip = true
		return
	}

	// Is category?
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		isCategory = true
		value = strings.Trim(line, "[]")
		return
	}

	// Is key/value pair?
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		skip = true
		return
	}

	key = parts[0]
	value = strings.TrimRight(parts[1], " ")

	return
}
