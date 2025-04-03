# geq

A Go port of [get-graphql-schema](https://github.com/prisma-labs/get-graphql-schema) - a command line utility for downloading a GraphQL schema from a GraphQL endpoint.

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