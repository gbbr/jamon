## jamon [![status](https://sourcegraph.com/api/repos/github.com/gbbr/jamon/.badges/status.png)](https://sourcegraph.com/github.com/gbbr/jamon)

Jamon is a delicious configuration file parser to be used with your application. A configuration file may look like this:

```ini
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
api.view=${base}/users/view
```

Categories are optional. You can have key/value pairs with no categories at the beginning of your configuration file. Comments can also be trailing and are allowed after key's values (ie. `key=value # Trailing comment`)

All keys and values are strings by default, if you need to convert to other types, use the amazing [strconv](http://golang.org/pkg/strconv/) package from Go's standard library.

#### Usage

To load a configuration file:

```go
config, err := jamon.Load("filename.config")

// For categorized keys:
config.Group("defaults").Get("user")

// For root-level keys:
config.Get("address")
```

Key & category getters do not return errors to allow chainability. If you specifically want to check whether a value exist boolean functions are provided, such as:

```go
// For root level
config.HasGroup("category_name")
config.Has("key_name")

// For categories
config.Group("defaults").Has("key_name")
```

[View me on GoDoc.org](http://godoc.org/github.com/gbbr/jamon)

#### Substitutions

Substitutions only support [alfanumeric values and dots](https://github.com/gbbr/jamon/blob/master/jamon.go#L61), so it is recommended that keys follow the same pattern. Substitutions are replaced in order of priority: first the group is checked, and next the root level. Cross-group substitutions are not allowed.

#### Notes

Hoisting not supported!  
Both keys and values are 100% case-sensitive
