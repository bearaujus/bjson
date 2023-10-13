# BJSON - JSON Library with Unique Features in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bearaujus/bjson/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bearaujus/bjson)](https://goreportcard.com/report/github.com/bearaujus/bjson)

BJSON is a Go package that takes JSON handling to the next level, offering additional features for flexibility and convenience. 
It empowers you to escape and unescape individual JSON elements and perform CRUD operations within JSON structures effortlessly.

## Key Usecase

- **Adaptive model**: BJSON is ideal for working with adaptive models, which are models that can be changed without having to change the code.

- **Escape and unescape elements**: BJSON provides functions for escaping and unescaping individual JSON elements. This can be useful for sanitizing data or for working with data that contains special characters.

- **CRUD**: BJSON provides functions for creating, reading, updating, and deleting JSON elements. This makes it easy to manipulate JSON data in Go.

- **IO**: BJSON provides function for unmarshal read and marshal write. This can be useful when your usecase is to handle many IO operations

## Installation

To install BJSON, you can run the following command:

```bash
go get github.com/bearaujus/bjson
```

## Usage

You can start the magic from here:

```go
package main

import (
	"github.com/bearaujus/bjson"
)

func main() {
	// from string
	bjson.NewBJSON(`{"name":"john","age":12}`)

	// from struct
	bjson.NewBJSON(struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{Name: "john", Age: 12})

	// from json object
	bjson.NewBJSON(map[string]interface{}{"name": "john", "age": 12})

	// from json array
	bjson.NewBJSON([]interface{}{"john", 12})
	
	// from file
	bjson.NewBJSONFromFile("path/to/file.json")

	// and more... checkout init.go for the detailed info
}

```


## TODO

- Refactor
- Add examples
- Add Iterator to iterate within the object
- Wrap errors and improve error messages
- Improve documentation

## License

This project is licensed under the MIT License - see the LICENSE file for details.
