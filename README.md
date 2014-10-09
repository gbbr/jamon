## jamon 

Jamon is a delicious configuration file parser to be used with your application. A configuration file may look like this:

```objectivec
address=127.0.0.1:123

[defaults]
user=gabriel
email=ga@stripetree.co.uk

# Comments use hash symbol
[category]
key=given_value
name.space.key=value with spaces
```

Categories are optional. You can have key/value pairs with no categories at the beginning of your configuration file. Comments can also be trailing and are allowed after key's values (ie. `key=value # Trailing comment`)

#### Usage

To load a configuration file:

```go
config := jamon.Load("filename.config")

// For categorized keys:
config.Category("my_category").Get("my_key")

// For non-categorized keys:
config.Get("my_key")
```

Key & category getters do not return errors to allow chainability. If you specifically want to check whether a value exist boolean functions are provided, such as:

```go
// For root level
config.HasCategory("category_name")
config.HasKey("key_name")

// For categories
config.Category("defaults").HasKey("key_name")
```
