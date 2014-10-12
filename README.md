## jamon 

Jamon is a delicious configuration file parser to be used with your application. A configuration file may look like this:

```objectivec
ip=127.0.0.1
port=123
address=${ip}:${port}
base=www.myaddr.com

[defaults]
user=gabriel
email=ga@stripetree.co.uk

# Comments use hash symbol
[routes]
api.users=${base}/users
api.view=/users/view
```

Categories are optional. You can have key/value pairs with no categories at the beginning of your configuration file. Comments can also be trailing and are allowed after key's values (ie. `key=value # Trailing comment`)

All keys and values are strings by default, if you need to convert to other types, use the amazing [strconv](http://golang.org/pkg/strconv/) package from Go's standard library.

#### Usage

To load a configuration file:

```go
config := jamon.Load("filename.config")

// For categorized keys:
config.Group("defaults").Key("user")

// For root-level keys:
config.Key("address")
```

Key & category getters do not return errors to allow chainability. If you specifically want to check whether a value exist boolean functions are provided, such as:

```go
// For root level
config.HasGroup("category_name")
config.HasKey("key_name")

// For categories
config.Group("defaults").HasKey("key_name")
```

#### Substitutions

Substitutions only support alfanumeric values and dots, so it is recommended that keys follow the same pattern. Substitutions are replaced in order of priority: first the group is checked, and next the root level. Cross-group substitutions are not allowed.

Hoisting not supported!
