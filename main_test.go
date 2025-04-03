package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTypeRefToString(t *testing.T) {
	tests := []struct {
		name     string
		typeRef  TypeRef
		expected string
	}{
		{
			name: "Scalar type",
			typeRef: TypeRef{
				Kind: "SCALAR",
				Name: "String",
			},
			expected: "String",
		},
		{
			name: "Non-null scalar type",
			typeRef: TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "SCALAR",
					Name: "String",
				},
			},
			expected: "String!",
		},
		{
			name: "List type",
			typeRef: TypeRef{
				Kind: "LIST",
				OfType: &TypeRef{
					Kind: "SCALAR",
					Name: "String",
				},
			},
			expected: "[String]",
		},
		{
			name: "Non-null list type",
			typeRef: TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Kind: "SCALAR",
						Name: "String",
					},
				},
			},
			expected: "[String]!",
		},
		{
			name: "List of non-null types",
			typeRef: TypeRef{
				Kind: "LIST",
				OfType: &TypeRef{
					Kind: "NON_NULL",
					OfType: &TypeRef{
						Kind: "SCALAR",
						Name: "String",
					},
				},
			},
			expected: "[String!]",
		},
		{
			name: "Non-null list of non-null types",
			typeRef: TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Kind: "NON_NULL",
						OfType: &TypeRef{
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
			result := typeRefToString(test.typeRef)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestGenerateSDL(t *testing.T) {
	// Sample minimal introspection response for testing
	responseJSON := `{
		"data": {
			"__schema": {
				"queryType": { "name": "Query" },
				"mutationType": { "name": "Mutation" },
				"subscriptionType": null,
				"types": [
					{
						"kind": "OBJECT",
						"name": "Query",
						"description": "The root query object",
						"fields": [
							{
								"name": "user",
								"description": "Get a user by ID",
								"args": [
									{
										"name": "id",
										"description": "The user ID",
										"type": {
											"kind": "NON_NULL",
											"name": null,
											"ofType": {
												"kind": "SCALAR",
												"name": "ID"
											}
										},
										"defaultValue": null
									}
								],
								"type": {
									"kind": "OBJECT",
									"name": "User"
								}
							}
						],
						"inputFields": null,
						"interfaces": [],
						"enumValues": null,
						"possibleTypes": null
					},
					{
						"kind": "OBJECT",
						"name": "User",
						"description": "A user in the system",
						"fields": [
							{
								"name": "id",
								"description": "The unique ID of the user",
								"args": [],
								"type": {
									"kind": "NON_NULL",
									"name": null,
									"ofType": {
										"kind": "SCALAR",
										"name": "ID"
									}
								}
							},
							{
								"name": "name",
								"description": "The name of the user",
								"args": [],
								"type": {
									"kind": "SCALAR",
									"name": "String"
								}
							}
						],
						"inputFields": null,
						"interfaces": [],
						"enumValues": null,
						"possibleTypes": null
					},
					{
						"kind": "SCALAR",
						"name": "ID",
						"description": "The ID scalar type",
						"fields": null,
						"inputFields": null,
						"interfaces": null,
						"enumValues": null,
						"possibleTypes": null
					},
					{
						"kind": "SCALAR",
						"name": "String",
						"description": "The String scalar type",
						"fields": null,
						"inputFields": null,
						"interfaces": null,
						"enumValues": null,
						"possibleTypes": null
					},
					{
						"kind": "OBJECT",
						"name": "Mutation",
						"description": "The root mutation object",
						"fields": [
							{
								"name": "createUser",
								"description": "Create a new user",
								"args": [
									{
										"name": "input",
										"description": "The user input",
										"type": {
											"kind": "NON_NULL",
											"name": null,
											"ofType": {
												"kind": "INPUT_OBJECT",
												"name": "CreateUserInput"
											}
										},
										"defaultValue": null
									}
								],
								"type": {
									"kind": "OBJECT",
									"name": "User"
								}
							}
						],
						"inputFields": null,
						"interfaces": [],
						"enumValues": null,
						"possibleTypes": null
					},
					{
						"kind": "INPUT_OBJECT",
						"name": "CreateUserInput",
						"description": "Input for creating a user",
						"fields": null,
						"inputFields": [
							{
								"name": "name",
								"description": "The name of the user",
								"type": {
									"kind": "NON_NULL",
									"name": null,
									"ofType": {
										"kind": "SCALAR",
										"name": "String"
									}
								},
								"defaultValue": null
							},
							{
								"name": "role",
								"description": "The role of the user",
								"type": {
									"kind": "ENUM",
									"name": "UserRole"
								},
								"defaultValue": "\"USER\""
							}
						],
						"interfaces": null,
						"enumValues": null,
						"possibleTypes": null
					},
					{
						"kind": "ENUM",
						"name": "UserRole",
						"description": "The role of a user",
						"fields": null,
						"inputFields": null,
						"interfaces": null,
						"enumValues": [
							{
								"name": "ADMIN",
								"description": "Administrator role"
							},
							{
								"name": "USER",
								"description": "Regular user role"
							}
						],
						"possibleTypes": null
					}
				],
				"directives": [
					{
						"name": "include",
						"description": "Directs the executor to include this field or fragment only when the argument is true.",
						"locations": ["FIELD", "FRAGMENT_SPREAD", "INLINE_FRAGMENT"],
						"args": [
							{
								"name": "if",
								"description": "Included when true.",
								"type": {
									"kind": "NON_NULL",
									"name": null,
									"ofType": {
										"kind": "SCALAR",
										"name": "Boolean"
									}
								},
								"defaultValue": null
							}
						]
					},
					{
						"name": "skip",
						"description": "Directs the executor to skip this field or fragment when the argument is true.",
						"locations": ["FIELD", "FRAGMENT_SPREAD", "INLINE_FRAGMENT"],
						"args": [
							{
								"name": "if",
								"description": "Skipped when true.",
								"type": {
									"kind": "NON_NULL",
									"name": null,
									"ofType": {
										"kind": "SCALAR",
										"name": "Boolean"
									}
								},
								"defaultValue": null
							}
						]
					}
				]
			}
		}
	}`

	var response IntrospectionResponse
	if err := json.Unmarshal([]byte(responseJSON), &response); err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	schema := generateSDL(response)

	// Check for expected elements in the SDL
	expectedElements := []string{
		"type Query {",
		"# The root query object",
		"user(",
		"# Get a user by ID",
		"# A user in the system",
		"type User {",
		"id: ID!",
		"name: String",
		"input CreateUserInput {",
		"# Input for creating a user",
		"name: String!",
		"role: UserRole = \"USER\"",
		"enum UserRole {",
		"# Administrator role",
		"ADMIN",
		"# Regular user role",
		"USER",
		"scalar ID",
		"scalar String",
		"directive @include(if: Boolean!)",
		"directive @skip(if: Boolean!)",
	}

	for _, expected := range expectedElements {
		if !contains(schema, expected) {
			t.Errorf("Expected SDL to contain %q, but it doesn't.\nActual: %s", expected, schema)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}