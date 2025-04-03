# geq

A command-line utility, written in Go, for downloading GraphQL schemas via introspection from a GraphQL endpoint. It serves as a Go alternative to tools like [get-graphql-schema](https://github.com/prisma-labs/get-graphql-schema).

**Why download a GraphQL schema?**

*   **Local Development:** Use the schema locally for enhanced tooling, such as autocompletion and validation in your IDE.
*   **Code Generation:** Generate client-side or server-side code (types, resolvers, SDKs) based on the schema structure.
*   **Documentation:** Keep an up-to-date reference of the API's structure.
*   **Schema Diffing:** Compare different versions of a schema to track changes over time.
*   **Offline Analysis:** Analyze the schema structure without needing constant access to the endpoint.

`geq` allows you to fetch the schema in either standard GraphQL Schema Definition Language (SDL) format or as the raw introspection JSON.

Additionally, the `--minify` option provides a description-free version of the schema (`schema.min.graphql` or `schema.min.json`). This is particularly useful for:

*   **Automated Tooling:** When tools only need the structural information (types, fields, arguments) without documentation comments.
*   **Reduced File Size:** Creating a more compact schema file when descriptions are unnecessary.
*   **LLM Context:** Providing schema structure to Large Language Models (LLMs) more efficiently by omitting descriptive text and reducing token count.

## Installation

```
go install github.com/pzurek/geq@latest
```

## Usage

```
geq --endpoint https://your-graphql-endpoint.com
```

### Options

- `-e`, `--endpoint`: The GraphQL endpoint URL (required)
- `-H`, `--header`: HTTP header in the format 'name: value'
    - Example for authentication: `--header "Authorization: YOUR_API_KEY"`
- `-o`, `--output`: Output file path for the schema (defaults to `schema.graphql` or `schema.json`)
- `-j`, `--json`: Output schema as JSON instead of SDL
- `-m`, `--minify`: Generate an additional minified schema file (no descriptions) named `schema.min.graphql` or `schema.min.json`.
- `-v`, `--version`: Show version information

## Development

Clone the repository and run:

```
go build
```