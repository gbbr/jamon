/*
Package Jamon is an INI-like configuration file parser. An example configuration
file may look like this:
	address=127.0.0.1:1234 # root-level values

	[defaults]
	key=value
	name=Gabriel

	[Group]
	key=value
Trailing comments are also allowed, and root-level keys are only accepted at the
top of the file
*/
package jamon

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Internal name for group that holds settings at root-level
const defaultGroup = "JAMON.ROOT_GROUP"

// A configuration type may hold multiple categories of settings
type Config map[string]Group

// A group holds key-value pairs of settings
type Group map[string]string

// Returns the value of a root-level key
func (c Config) Key(key string) string { return c[defaultGroup].Key(key) }

// Verifies the existence of a root-level key
func (c Config) HasKey(key string) bool {
	_, ok := c[defaultGroup][key]
	return ok
}

// Returns a group by name. If the group does not exist, an empty category
// is returned. This is to avoid multiple return values in order to facilitate
// chaining.
func (c Config) Group(name string) Group { return c[name] }

// Verifies if a category exists
func (c Config) HasGroup(category string) bool {
	_, ok := c[category]
	return ok
}

// Returns a key from a category
func (c Group) Key(key string) string { return c[key] }

// Verifies if the category has a key
func (c Group) HasKey(key string) bool {
	_, ok := c[key]
	return ok
}

// Regexp for substitions
var regexSubst = regexp.MustCompile(`\$\{([\w\s.]*)\}`)

// Loads a configuration file
func Load(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := Config{}
	currentGroup := defaultGroup

	for scanner.Scan() {
		isGroup, value, key, skip := parseLine(scanner.Text())

		switch {
		case skip:
			continue

		case isGroup:
			currentGroup = value
			continue

		case config[currentGroup] == nil:
			config[currentGroup] = make(Group)
			fallthrough

		default:
			compile := func(r string) string {
				k := r[2 : len(r)-1]

				if _, ok := config[currentGroup][k]; ok {
					return config[currentGroup][k]
				}

				if _, ok := config[defaultGroup][k]; ok {
					return config[defaultGroup][k]
				}

				return r
			}

			val := regexSubst.ReplaceAllStringFunc(value, compile)
			config[currentGroup][key] = val
		}
	}

	return config, nil
}

// Attempts to parse an entry in the config file. The first return value specifies
// whether 'value' is the name of a category or the value of a key. Skip indicates
// whether the line was a comment or could not be parsed.
func parseLine(line string) (isGroup bool, value, key string, skip bool) {
	line = strings.SplitN(line, "#", 2)[0]
	line = strings.Trim(line, " \t\r")

	// Is comment or empty line?
	if len(line) == 0 {
		skip = true
		return
	}

	// Is category?
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		isGroup = true
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
