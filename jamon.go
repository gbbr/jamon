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

// Internal name for category that holds settings without one
const defaultCategory = "JAMON.NO_CATEGORY"

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
// whether 'value' is the name of a category or the value of a key.
func parseLine(line string) (isCategory bool, value, key string, skip bool) {
	line = strings.SplitN(line, "#", 2)[0]
	line = strings.Trim(line, " \t\r")

	// Is comment?
	if strings.HasPrefix(line, "#") || len(line) == 0 {
		skip = true
		return
	}

	// Is category?
	if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
		isCategory = true
		value = strings.Trim(line, "[]")
		return
	}

	// Attempt to parse key/value pair
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		skip = true
		return
	}

	// Trim end-of-line comments
	key = parts[0]
	value = strings.TrimRight(parts[1], " ")

	return
}

// Returns the value of a key that is not in any category. These keys should
// be placed at the top of the file with no title if desired.
func (c Config) Get(key string) string {
	category, ok := c["JAMON.NO_CATEGORY"]
	if !ok {
		return ""
	}

	return category.Get(key)
}

// Returns a category by name. If the category does not exist, an empty category
// is returned. Errors are not returned here in order to allow chaining.
func (c Config) Category(name string) Category { return c[name] }

// Returns a key from a category
func (c Category) Get(key string) string { return c[key] }
