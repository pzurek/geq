# geq

A Go port of [get-graphql-schema](https://github.com/prisma-labs/get-graphql-schema) - a command line utility to download a GraphQL schema from a GraphQL endpoint.

## Installation

```
go install github.com/pzurek/geq@latest
```

## Usage

```
geq --endpoint https://your-graphql-endpoint.com
```

### Options

- `--endpoint`: The GraphQL endpoint URL (required)
- `--header`: HTTP header in the format 'name: value'
    - Example for authentication: `--header "Authorization: YOUR_API_KEY"`
- `--json`: Output schema as JSON instead of SDL
- `--version`: Show version information
- `--minify`: Generate an additional minified schema file (no descriptions) named `schema.min.graphql` or `schema.min.json`.

## Development

Clone the repository and run:

```