# BJSON
BJSON is a Go package that provides a JSON library with additional features for flexibility, including the ability to 
escape and unescape individual JSON elements and the ability to set a root JSON element for marshaling. 
With BJSON, you don't need to declare models beforehand to unmarshal JSON data.

## Installation
Use `go get` to install BJSON:
```bash
go get -u github.com/bearaujus/bjson
```

## Usage
### Importing the package
```go
import "github.com/bearaujus/bjson"
```

### Creating a new BJSON object
```go
jsonData := []byte(`{"foo": "bar", "num": 42, "nested": {"a": [1, 2, 3], "b": true}}`)
bj, err := bjson.NewBJSON(jsonData)
if err != nil {
    log.Fatal(err)
}
```

### Marshaling a BJSON object
```go
jsonData, err := bj.MarshalJSON()
if err != nil {
    log.Fatal(err)
}
```

### Marshaling a BJSON object with indentation
```go
jsonData, err := bj.MarshalJSONPretty()
if err != nil {
    log.Fatal(err)
}
```

### Setting a root JSON element for marshaling
```go
err := bj.SetMarshalRootJSONElement([]string{"nested", "a"})
if err != nil {
    log.Fatal(err)
}
```

### Resetting the root JSON element
```go
bj.ResetMarshalRootJSONElement()
```

### Removing an element
```go
err := bj.RemoveElement([]string{"nested", "a"})
if err != nil {
    log.Fatal(err)
}
```

### Escaping an element
```go
err := bj.EscapeJSONElement([]string{"nested", "a"})
if err != nil {
    log.Fatal(err)
}
```

### Unescaping an element
```go
err := bj.UnescapeJSONElement([]string{"nested", "a"})
if err != nil {
    log.Fatal(err)
}
```

## TODO
- Add Unittests

## License
This package is licensed under the MIT license. See the LICENSE file for more details.
