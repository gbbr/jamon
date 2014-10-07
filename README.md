_WORK IN PROGRESS_


## jamon 

Jamon is a delicious configuration file parser to be used with your application. A configuration file may look like this:

```
[category]
key=value
key2=value2

# This category is important
[next.category]
long.key=other value # Trailing comment
name.space.key=value
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
