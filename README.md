## jamon 

Jamon is a delicious configuration file parser to be used with your application. A configuration file may look like this:

```vbnet
address=127.0.0.1:123

[defaults]
user=gabriel
email=ga@stripetree.co.uk

# This category is important
[category]
long.key=given_value
name.space.key=value with spaces # Trailing comment
```

Categories are optional. You can have key/value pairs with no categories at the beginning of your configuration file.

#### Usage

To load a configuration file:

```go
config := jamon.Load("filename.config")

// For categorized keys:
config.Category("my_category").Get("my_key")

// For non-categorized keys:
config.Get("my_key")
```
