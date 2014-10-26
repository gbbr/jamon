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
top of the file.
*/
package jamon

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Internal name for root-level group.
const rootGroup = "JAMON.ROOT_GROUP"

// A configuration holds keys and/or groups of keys.
type Config map[string]Group

// A group holds key-value pairs.
type Group map[string]string

// Returns the value of a root-level key.
func (c Config) Get(key string) string { return c[rootGroup].Get(key) }

// Verifies the existence of a root-level key.
func (c Config) Has(key string) bool {
	_, ok := c[rootGroup][key]
	return ok
}

// Returns a group in the configuration file or an empty one if it doesn't exist.
func (c Config) Group(name string) Group { return c[name] }

// Verifies if a group exists.
func (c Config) HasGroup(category string) bool {
	_, ok := c[category]
	return ok
}

// Returns a key from the group.
func (c Group) Get(key string) string { return c[key] }

// Verifies if the group contains the key.
func (c Group) Has(key string) bool {
	_, ok := c[key]
	return ok
}

// Regexp for substitions
var regexSubst = regexp.MustCompile(`\$\{([^\}\{]*)\}`)

// Loads a configuration file.
func Load(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	grp := rootGroup
	cfg := Config{}

	for scanner.Scan() {
		isGroup, val, key, skip := parseLine(scanner.Text())

		switch {
		case skip:
			continue

		case isGroup:
			grp = val
			continue

		case cfg[grp] == nil:
			cfg[grp] = make(Group)
			fallthrough

		default:
			replaceFn := func(r string) string {
				sub := r[2 : len(r)-1]

				// Is replacement in own group?
				if _, ok := cfg[grp][sub]; ok {
					return cfg[grp][sub]
				}
				// Is replacement in root group?
				if _, ok := cfg[rootGroup][sub]; ok {
					return cfg[rootGroup][sub]
				}
				// If it's not found, no change happens
				return r
			}

			cfg[grp][key] = regexSubst.ReplaceAllStringFunc(val, replaceFn)
		}
	}

	return cfg, nil
}

// Attempts to parse an entry in the config file. The first return value specifies
// whether 'value' is the name of a category or the value of a key. Skip indicates
// whether the line was a comment or could not be parsed.
func parseLine(line string) (isGroup bool, val, key string, skip bool) {
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
		val = strings.Trim(line, "[]")
		return
	}
	// Is key/value pair?
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		skip = true
		return
	}

	key, val = parts[0], strings.TrimRight(parts[1], " ")
	return
}
