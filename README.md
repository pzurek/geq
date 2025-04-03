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
- `--json`: Output schema as JSON instead of SDL

## Development

Clone the repository and run:

```
go build
```
