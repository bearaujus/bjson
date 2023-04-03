package main

import (
	"fmt"
	"log"

	"github.com/bearaujus/bjson"
)

func main() {
	// JSON string to be used in this example
	jsonString := `{"name":"John Doe","age":30,"data":{"address":"{\"city\":\"New York\",\"country\":\"USA\"}", "phone": {"foo": "bar", "x": "y"}}}`

	// Create a BJSON object from a JSON string
	bj, err := bjson.NewBJSONFromString(jsonString)
	if err != nil {
		log.Fatal(err)
	}

	// Unescape a JSON element
	err = bj.UnescapeJSONElement([]string{"data", "address"})
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the BJSON object into a formatted JSON string after unescaping an element
	data, err := bj.MarshalJSONPretty()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Marshalled JSON after unescaping an element:\n", string(data))

	// Escape a JSON element
	err = bj.EscapeJSONElement([]string{"data", "phone"})
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the BJSON object into a formatted JSON string after escaping an element
	data, err = bj.MarshalJSONPretty()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Marshalled JSON after escaping an element:\n", string(data))

	// Remove a JSON element
	err = bj.RemoveElement([]string{"data", "phone"})
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the BJSON object into a formatted JSON string after removing an element
	data, err = bj.MarshalJSONPretty()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Marshalled JSON after removing an element:\n", string(data))

	// Set the root JSON element for marshaling
	err = bj.SetMarshalRootJSONElement([]string{"data"})
	if err != nil {
		log.Fatal(err)
	}

	// Marshal the BJSON object into a formatted JSON string
	data, err = bj.MarshalJSONPretty()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Marshalled JSON after set marshall root element:\n", string(data))
}
