# BJSON - A Model-Less JSON Library in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/bearaujus/bjson/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/bearaujus/bjson)](https://goreportcard.com/report/github.com/bearaujus/bjson)

BJSON is a Go library that provides a simple way to work with JSON data. It is designed to be flexible and easy to use, and it does not require you to declare any models. This makes it ideal for working with dynamic JSON data, such as the responses from APIs.

## Features

- **Adaptive model**: BJSON is ideal for working with adaptive models, which are models that can be changed without having to change the code.
- **Escape and unescape elements**: BJSON provides functions for escaping and unescaping individual JSON elements. This can be useful for sanitizing data or for working with data that contains special characters.
- **CRUD**: BJSON provides functions for creating, reading, updating, and deleting JSON elements. This makes it easy to manipulate JSON data in Go.
- **IO**: BJSON provides function for unmarshal read and marshal write. This can be useful when your usecase is to handle many IO operations

## Installation

To install BJSON, you can run the following command:

```shell
go get github.com/bearaujus/bjson
```

## Usage

The following example shows how to use BJSON to process a simple JSON:

```go
import (
    "github.com/bearaujus/bjson"
)

func main() {
    // Create a new BJSON object from a JSON string.
    jsonStr := `{ "name": "John Doe", "age": 30 }`
    bj, err := bjson.NewBJSON(jsonStr)
    if err != nil {
        // Handle error
    }

    // Add a new element to the JSON object.
    err = bj.AddElement("occupation", "Software Engineer")
    if err != nil {
        // Handle error
    }

    // Get the value of the "occupation" element.
    occupation, err := bj.GetElement("occupation")
    if err != nil {
        // Handle error
    }

    // Set the value of the "age" element.
    err = bj.SetElement(25, "age")
    if err != nil {
        // Handle error
    }

    // Print the JSON object.
    fmt.Println(bj.String())
}

```
```json
{"name":"John Doe","age":25,"occupation":"Software Engineer"}
```

## TODO
- Add documentation to the code
- Add more examples to the README.md
- Improve error handling
- Refactor the code to be simpler

## License

This project is licensed under the MIT License - see the LICENSE file for details.
