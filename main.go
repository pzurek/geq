package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/pzurek/geq/pkg/geq"
)

// Build time variables
var version = "dev"
var buildTime = "unset"

// writeSchemaFile handles file writing and console output for the CLI.
func writeSchemaFile(outputPath string, content string) error {
	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing schema to file '%s': %v\n", outputPath, err)
		return err // Return error for main to handle exit
	}
	fmt.Printf("Schema successfully saved to %s\n", outputPath)
	return nil
}

func main() {
	// Parse command line arguments
	endpoint := flag.String("endpoint", "", "The GraphQL endpoint URL")
	header := flag.String("header", "", "Header in the format 'name: value'")
	outputFile := flag.String("output", "", "Output file path for the schema (SDL or JSON)")
	asJSON := flag.Bool("json", false, "Output as JSON")
	versionFlag := flag.Bool("version", false, "Show version information")
	minify := flag.Bool("minify", false, "Generate an additional minified schema file (no descriptions)")

	// Short flag aliases
	flag.StringVar(endpoint, "e", *endpoint, "The GraphQL endpoint URL (shorthand)")
	flag.StringVar(header, "H", *header, "Header in the format 'name: value' (shorthand)")
	flag.StringVar(outputFile, "o", *outputFile, "Output file path (shorthand)")
	flag.BoolVar(asJSON, "j", *asJSON, "Output as JSON (shorthand)")
	flag.BoolVar(versionFlag, "v", *versionFlag, "Show version information (shorthand)")
	flag.BoolVar(minify, "m", *minify, "Generate minified schema (shorthand)")

	flag.Parse()

	// Show version if requested
	if *versionFlag {
		fmt.Printf("geq version %s (built %s)\n", version, buildTime)
		return
	}

	// Check if URL is provided
	if *endpoint == "" {
		fmt.Println("Error: GraphQL endpoint URL is required")
		flag.Usage()
		os.Exit(1)
	}

	// Fetch schema data using the library function
	introspectionJSON, err := geq.FetchIntrospectionJSON(*endpoint, *header)
	if err != nil {
		fmt.Printf("Error fetching schema data: %v\n", err)
		os.Exit(1)
	}

	// Determine main output path and format
	mainOutputPath := *outputFile
	outputIsJSON := *asJSON
	mainSchemaContent := ""

	if mainOutputPath == "" {
		if outputIsJSON {
			mainOutputPath = "schema.json"
		} else {
			mainOutputPath = "schema.graphql"
		}
	}

	// Generate main schema content (SDL or JSON)
	if outputIsJSON {
		// JSON output logic
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(introspectionJSON), "", "  "); err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
			os.Exit(1)
		}
		mainSchemaContent = prettyJSON.String()
	} else {
		// SDL output logic
		var introspectionResp geq.IntrospectionResponse
		if err := json.Unmarshal([]byte(introspectionJSON), &introspectionResp); err != nil {
			// Provide more context on JSON parsing error
			fmt.Printf("Error parsing introspection JSON response: %v\n", err)
			// Optionally print a snippet of the received JSON if it's not too large
			snippet := introspectionJSON
			if len(snippet) > 200 {
				snippet = snippet[:200] + "..."
			}
			fmt.Printf("Received JSON snippet: %s\n", snippet)
			os.Exit(1)
		}
		mainSchemaContent = geq.GenerateSDL(introspectionResp)
	}

	// Write main schema file using the local function
	err = writeSchemaFile(mainOutputPath, mainSchemaContent)
	if err != nil {
		os.Exit(1) // Exit if writing failed
	}

	// Generate and write minified schema if requested
	if *minify {
		minifiedOutputPath := ""
		minifiedSchemaContent := ""

		// Determine minified output path
		if outputIsJSON {
			minifiedOutputPath = "schema.min.json"
		} else {
			minifiedOutputPath = "schema.min.graphql"
		}

		// Generate minified schema content
		if outputIsJSON {
			// Minified JSON logic
			var compactJSON bytes.Buffer
			// Use json.Compact instead of Marshal for minification
			if err := json.Compact(&compactJSON, []byte(introspectionJSON)); err != nil {
				fmt.Printf("Error compacting JSON: %v\n", err)
				os.Exit(1)
			}
			minifiedSchemaContent = compactJSON.String()
		} else {
			// Minified SDL logic
			var introspectionResp geq.IntrospectionResponse
			if err := json.Unmarshal([]byte(introspectionJSON), &introspectionResp); err != nil {
				fmt.Printf("Error parsing introspection response for minify: %v\n", err)
				os.Exit(1)
			}
			minifiedSchemaContent = geq.GenerateMinifiedSDL(introspectionResp)
		}

		// Write minified schema file using the local function
		err = writeSchemaFile(minifiedOutputPath, minifiedSchemaContent)
		if err != nil {
			os.Exit(1) // Exit if writing failed
		}
	}
}
