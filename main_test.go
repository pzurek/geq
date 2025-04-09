package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/pzurek/geq/pkg/geq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "Update golden files")

func TestTypeRefToString(t *testing.T) {
	tests := []struct {
		name     string
		typeRef  geq.TypeRef
		expected string
	}{
		{
			name: "Scalar type",
			typeRef: geq.TypeRef{
				Kind: "SCALAR",
				Name: "String",
			},
			expected: "String",
		},
		{
			name: "Non-null scalar type",
			typeRef: geq.TypeRef{
				Kind: "NON_NULL",
				OfType: &geq.TypeRef{
					Kind: "SCALAR",
					Name: "String",
				},
			},
			expected: "String!",
		},
		{
			name: "List type",
			typeRef: geq.TypeRef{
				Kind: "LIST",
				OfType: &geq.TypeRef{
					Kind: "SCALAR",
					Name: "String",
				},
			},
			expected: "[String]",
		},
		{
			name: "Non-null list type",
			typeRef: geq.TypeRef{
				Kind: "NON_NULL",
				OfType: &geq.TypeRef{
					Kind: "LIST",
					OfType: &geq.TypeRef{
						Kind: "SCALAR",
						Name: "String",
					},
				},
			},
			expected: "[String]!",
		},
		{
			name: "List of non-null types",
			typeRef: geq.TypeRef{
				Kind: "LIST",
				OfType: &geq.TypeRef{
					Kind: "NON_NULL",
					OfType: &geq.TypeRef{
						Kind: "SCALAR",
						Name: "String",
					},
				},
			},
			expected: "[String!]",
		},
		{
			name: "Non-null list of non-null types",
			typeRef: geq.TypeRef{
				Kind: "NON_NULL",
				OfType: &geq.TypeRef{
					Kind: "LIST",
					OfType: &geq.TypeRef{
						Kind: "NON_NULL",
						OfType: &geq.TypeRef{
							Kind: "SCALAR",
							Name: "String",
						},
					},
				},
			},
			expected: "[String!]!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := geq.TypeRefToString(test.typeRef)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestGenerateSDL(t *testing.T) {
	// Read input introspection JSON
	inputJSONBytes, err := os.ReadFile(filepath.Join("testdata", "sample_introspection.json"))
	require.NoError(t, err, "Failed to read input JSON file")

	var response geq.IntrospectionResponse
	err = json.Unmarshal(inputJSONBytes, &response)
	require.NoError(t, err, "Failed to parse test JSON")

	// Generate actual SDL
	actualSDL := geq.GenerateSDL(response)

	// Define golden file path
	goldenFilePath := filepath.Join("testdata", "sample_schema.graphql")

	// Update golden file if requested
	if *update {
		err = os.WriteFile(goldenFilePath, []byte(actualSDL), 0644)
		require.NoError(t, err, "Failed to write golden file")
		t.Logf("Golden file updated: %s", goldenFilePath)
		return // Don't compare if we just updated
	}

	// Read expected SDL from golden file
	expectedSDLBytes, err := os.ReadFile(goldenFilePath)
	// Handle case where golden file doesn't exist yet
	if os.IsNotExist(err) {
		// Create the golden file with the current output for the first run
		err = os.WriteFile(goldenFilePath, []byte(actualSDL), 0644)
		require.NoError(t, err, "Failed to write initial golden file")
		t.Fatalf("Golden file %s did not exist. Created it with current output. Re-run tests.", goldenFilePath)
	} else {
		require.NoError(t, err, "Failed to read golden file")
	}

	// Compare actual vs expected
	assert.Equal(t, string(expectedSDLBytes), actualSDL, "Generated SDL does not match golden file %s", goldenFilePath)
}

func TestGenerateMinifiedSDL(t *testing.T) {
	// Read input introspection JSON
	inputJSONBytes, err := os.ReadFile(filepath.Join("testdata", "sample_introspection.json"))
	require.NoError(t, err, "Failed to read input JSON file")

	var response geq.IntrospectionResponse
	err = json.Unmarshal(inputJSONBytes, &response)
	require.NoError(t, err, "Failed to parse test JSON")

	// Generate actual minified SDL
	actualMinSDL := geq.GenerateMinifiedSDL(response)

	// Define golden file path
	goldenFilePath := filepath.Join("testdata", "sample_schema.min.graphql")

	// Update golden file if requested
	if *update {
		err = os.WriteFile(goldenFilePath, []byte(actualMinSDL), 0644)
		require.NoError(t, err, "Failed to write minified golden file")
		t.Logf("Minified golden file updated: %s", goldenFilePath)
		return // Don't compare if we just updated
	}

	// Read expected minified SDL from golden file
	expectedMinSDLBytes, err := os.ReadFile(goldenFilePath)
	// Handle case where golden file doesn't exist yet
	if os.IsNotExist(err) {
		// Create the golden file with the current output for the first run
		err = os.WriteFile(goldenFilePath, []byte(actualMinSDL), 0644)
		require.NoError(t, err, "Failed to write initial minified golden file")
		t.Fatalf("Minified golden file %s did not exist. Created it with current output. Re-run tests.", goldenFilePath)
	} else {
		require.NoError(t, err, "Failed to read minified golden file")
	}

	// Compare actual vs expected
	assert.Equal(t, string(expectedMinSDLBytes), actualMinSDL, "Generated minified SDL does not match golden file %s", goldenFilePath)
}
