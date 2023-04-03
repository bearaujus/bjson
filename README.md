# BJSON

BJSON is a Go package that provides a JSON library with additional features for flexibility, including the ability to
escape and unescape individual JSON elements and the ability to set a root JSON element for marshaling.
With BJSON, you don't need to declare models beforehand to unmarshal JSON data.

## Features:

- Read JSON data from a file, string, or byte slice
- Unmarshal and Marshal JSON data
- Set a root JSON element for marshaling
- Escape and Unescape JSON elements
- Remove JSON elements from the object
- Write marshaled JSON data to a file

## Installation

To install bjson, run the following command:

```bash
go get github.com/your-username/bjson
```

## Usage

To use bjson, you need to import the package:

```go
import "github.com/your-username/bjson"
```

## Creating a BJSON object

You can create a new BJSON object by unmarshaling JSON data from a file, string, or byte slice:

```go
// Unmarshal JSON data from a file
bj, err := bjson.NewBJSONFromFile("path/to/file.json")

// Unmarshal JSON data from a string
bj, err := bjson.NewBJSONFromString(`{"name": "John", "age": 30}`)

// Unmarshal JSON data from a byte slice
bj, err := bjson.NewBJSONFromByte([]byte(`{"name": "John", "age": 30}`))
```

## Marshaling and unmarshaling JSON data

You can marshal a BJSON object into a JSON string:

```go
data, err := bj.MarshalJSON()
```

You can also marshal a BJSON object into a formatted JSON string:

```go
data, err := bj.MarshalJSONPretty()
```

And you can unmarshal a JSON string into a BJSON object:

```go
err := bj.UnmarshalJSON([]byte(`{"name": "John", "age": 30}`))
```

## Manipulating JSON data

You can set the root JSON element for marshaling:

```go
err := bj.SetMarshalRootJSONElement([]string{"person", "address"})
```

You can reset the root JSON element:

```go
err := bj.ResetMarshalRootJSONElement()
```

You can remove a JSON element from the object:

```go
err := bj.RemoveElement([]string{"person", "address", "city"})
```

You can escape a JSON element by marshaling it into a JSON string:

```go
err := bj.EscapeJSONElement([]string{"person", "address"})
```

And you can unescape a JSON element by unmarshaling it from a JSON string:

```go
err := bj.UnescapeJSONElement([]string{"person", "address"})
```

## Writing JSON data to a file

You can write marshaled JSON data to a file:

```go
err := bj.WriteMarshalJSON("path/to/file.json", false)
```

By default, the JSON data is not formatted. If you want to format the JSON data, set the isPretty parameter to true.

## TODO

- Add more examples and use cases to the documentation.
- Improve error handling and error messages.
- Add support for get value from existing JSON elements.
- Add support for modifying existing JSON elements.
- Add support for merging two BJSON objects.
- Add support for converting BJSON to XML.
- Add support for converting BJSON to YAML.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
