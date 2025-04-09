# geq

A Go library and command-line utility for downloading GraphQL schemas via introspection from a GraphQL endpoint. It serves as a Go alternative to tools like [get-graphql-schema](https://github.com/prisma-labs/get-graphql-schema).

**Why download a GraphQL schema?**

*   **Local Development:** Use the schema locally for enhanced tooling, such as autocompletion and validation in your IDE.
*   **Code Generation:** Generate client-side or server-side code (types, resolvers, SDKs) based on the schema structure.
*   **Documentation:** Keep an up-to-date reference of the API's structure.
*   **Schema Diffing:** Compare different versions of a schema to track changes over time.
*   **Offline Analysis:** Analyze the schema structure without needing constant access to the endpoint.

`geq` allows you to fetch the schema in either standard GraphQL Schema Definition Language (SDL) format or as the raw introspection JSON.

Additionally, the minification feature provides a description-free version of the schema. This is particularly useful for:

*   **Automated Tooling:** When tools only need the structural information (types, fields, arguments) without documentation comments.
*   **Reduced File Size:** Creating a more compact schema file when descriptions are unnecessary.
*   **LLM Context:** Providing schema structure to Large Language Models (LLMs) more efficiently by omitting descriptive text and reducing token count.

## Installation

### CLI Tool

```/dev/null/install.sh#L1-2
go install github.com/pzurek/geq@latest
```

### Library

```/dev/null/go-get.sh#L1-2
go get github.com/pzurek/geq
```

## Usage

### CLI Usage

```/dev/null/cli-usage.sh#L1-2
geq --endpoint https://your-graphql-endpoint.com
```

#### CLI Options

- `-e`, `--endpoint`: The GraphQL endpoint URL (required)
- `-H`, `--header`: HTTP header in the format 'name: value'
    - Example for authentication: `--header "Authorization: YOUR_API_KEY"`
- `-o`, `--output`: Output file path for the schema (defaults to `schema.graphql` or `schema.json`)
- `-j`, `--json`: Output schema as JSON instead of SDL
- `-m`, `--minify`: Generate an additional minified schema file (no descriptions) named `schema.min.graphql` or `schema.min.json`.
- `-v`, `--version`: Show version information

### Library Usage

The `geq` library provides functions to fetch and process GraphQL schemas programmatically:

```/dev/null/library-example.go#L1-30
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pzurek/geq/pkg/geq"
)

func main() {
	// Fetch schema data using the library function
	endpoint := "https://your-graphql-endpoint.com"
	header := "Authorization: YOUR_API_KEY" // Optional

	introspectionJSON, err := geq.FetchIntrospectionJSON(endpoint, header)
	if err != nil {
		fmt.Printf("Error fetching schema data: %v\n", err)
		os.Exit(1)
	}

	// Parse introspection response
	var response geq.IntrospectionResponse
	err = json.Unmarshal([]byte(introspectionJSON), &response)
	if err != nil {
		fmt.Printf("Error parsing introspection response: %v\n", err)
		os.Exit(1)
	}

	// Generate SDL from introspection data
	sdl := geq.GenerateSDL(response)
	fmt.Println("Schema SDL:", sdl)

	// Generate minified SDL
	minifiedSDL := geq.GenerateMinifiedSDL(response)
	fmt.Println("Minified Schema SDL:", minifiedSDL)
}
```

### Key Library Functions

- `FetchIntrospectionJSON(endpoint, header string) (string, error)`: Fetches the raw introspection JSON from a GraphQL endpoint
- `GenerateSDL(response IntrospectionResponse) string`: Converts introspection response to SDL format
- `GenerateMinifiedSDL(response IntrospectionResponse) string`: Generates minified SDL without descriptions
- `TypeRefToString(typeRef TypeRef) string`: Utility function to convert type references to string representation

## Development

Clone the repository and run:

```/dev/null/build.sh#L1-2
go build
```

To run tests:

```/dev/null/tests.sh#L1-2
go test ./...
```