package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Version information - will be injected during build
var (
	version   = "dev"
	buildTime = "unknown"
)

// IntrospectionResponse represents the GraphQL introspection query response
type IntrospectionResponse struct {
	Data struct {
		Schema struct {
			QueryType struct {
				Name string `json:"name"`
			} `json:"queryType"`
			MutationType struct {
				Name string `json:"name"`
			} `json:"mutationType"`
			SubscriptionType struct {
				Name string `json:"name"`
			} `json:"subscriptionType"`
			Types []struct {
				Kind        string `json:"kind"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Fields      []struct {
					Name              string `json:"name"`
					Description       string `json:"description"`
					Args              []struct {
						Name         string `json:"name"`
						Description  string `json:"description"`
						Type         TypeRef `json:"type"`
						DefaultValue string  `json:"defaultValue"`
					} `json:"args"`
					Type TypeRef `json:"type"`
				} `json:"fields"`
				InputFields []struct {
					Name         string `json:"name"`
					Description  string `json:"description"`
					Type         TypeRef `json:"type"`
					DefaultValue string  `json:"defaultValue"`
				} `json:"inputFields"`
				Interfaces []TypeRef `json:"interfaces"`
				EnumValues []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"enumValues"`
				PossibleTypes []TypeRef `json:"possibleTypes"`
			} `json:"types"`
			Directives []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Locations   []string `json:"locations"`
				Args        []struct {
					Name         string `json:"name"`
					Description  string `json:"description"`
					Type         TypeRef `json:"type"`
					DefaultValue string  `json:"defaultValue"`
				} `json:"args"`
			} `json:"directives"`
		} `json:"__schema"`
	} `json:"data"`
}

// TypeRef represents a GraphQL type reference
type TypeRef struct {
	Kind   string  `json:"kind"`
	Name   string  `json:"name"`
	OfType *TypeRef `json:"ofType"`
}

func main() {
	// Parse command line arguments
	endpoint := flag.String("endpoint", "", "The GraphQL endpoint URL")
	header := flag.String("header", "", "Header in the format 'name: value'")
	json := flag.Bool("json", false, "Output as JSON")
	versionFlag := flag.Bool("version", false, "Show version information")

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

	// Fetch schema
	schema, err := fetchSchema(*endpoint, *header, *json)
	if err != nil {
		fmt.Printf("Error fetching schema: %v\n", err)
		os.Exit(1)
	}

	// Output schema
	fmt.Println(schema)
}

func fetchSchema(endpoint, headerStr string, asJSON bool) (string, error) {
	// Full introspection query - this is the standard query for fetching the complete schema
	introspectionQuery := `
	query IntrospectionQuery {
		__schema {
			queryType { name }
			mutationType { name }
			subscriptionType { name }
			types {
				kind
				name
				description
				fields(includeDeprecated: true) {
					name
					description
					args {
						name
						description
						type {
							kind
							name
							ofType {
								kind
								name
								ofType {
									kind
									name
									ofType {
										kind
										name
										ofType {
											kind
											name
											ofType {
												kind
												name
												ofType {
													kind
													name
													ofType {
														kind
														name
													}
												}
											}
										}
									}
								}
							}
						}
						defaultValue
					}
					type {
						kind
						name
						ofType {
							kind
							name
							ofType {
								kind
								name
								ofType {
									kind
									name
									ofType {
										kind
										name
										ofType {
											kind
											name
											ofType {
												kind
												name
												ofType {
													kind
													name
												}
											}
										}
									}
								}
							}
						}
					}
				}
				inputFields {
					name
					description
					type {
						kind
						name
						ofType {
							kind
							name
							ofType {
								kind
								name
								ofType {
									kind
									name
									ofType {
										kind
										name
										ofType {
											kind
											name
											ofType {
												kind
												name
												ofType {
													kind
													name
												}
											}
										}
									}
								}
							}
						}
					}
					defaultValue
				}
				interfaces {
					kind
					name
					ofType {
						kind
						name
						ofType {
							kind
							name
							ofType {
								kind
								name
								ofType {
									kind
									name
									ofType {
										kind
										name
										ofType {
											kind
											name
											ofType {
												kind
												name
											}
										}
									}
								}
							}
						}
					}
				}
				enumValues(includeDeprecated: true) {
					name
					description
				}
				possibleTypes {
					kind
					name
					ofType {
						kind
						name
						ofType {
							kind
							name
							ofType {
								kind
								name
								ofType {
									kind
									name
									ofType {
										kind
										name
										ofType {
											kind
											name
											ofType {
												kind
												name
											}
										}
									}
								}
							}
						}
					}
				}
			}
			directives {
				name
				description
				locations
				args {
					name
					description
					type {
						kind
						name
						ofType {
							kind
							name
							ofType {
								kind
								name
								ofType {
									kind
									name
									ofType {
										kind
										name
										ofType {
											kind
											name
											ofType {
												kind
												name
												ofType {
													kind
													name
												}
											}
										}
									}
								}
							}
						}
					}
					defaultValue
				}
			}
		}
	}`

	// Prepare the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query": introspectionQuery,
	})
	if err != nil {
		return "", fmt.Errorf("error creating request body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Add custom header if provided
	if headerStr != "" {
		parts := strings.SplitN(headerStr, ":", 2)
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			req.Header.Set(name, value)
		} else {
			return "", fmt.Errorf("invalid header format. Expected 'name: value', got '%s'", headerStr)
		}
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	// Check if status code is successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned error: %s", body)
	}

	// If asJSON flag is true, just return the raw JSON
	if asJSON {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
			return "", fmt.Errorf("error formatting JSON: %w", err)
		}
		return prettyJSON.String(), nil
	}

	// Parse response to extract schema
	var introspectionResp IntrospectionResponse
	if err := json.Unmarshal(body, &introspectionResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	// Convert introspection response to SDL format
	return generateSDL(introspectionResp), nil
}

// generateSDL converts the introspection response to SDL (Schema Definition Language) format
func generateSDL(response IntrospectionResponse) string {
	var sb strings.Builder

	// Track types that have been printed to avoid duplicates
	printedTypes := make(map[string]bool)

	// Process schema types
	for _, typeObj := range response.Data.Schema.Types {
		// Skip internal GraphQL types that start with "__"
		if strings.HasPrefix(typeObj.Name, "__") {
			continue
		}

		if printedTypes[typeObj.Name] {
			continue
		}
		printedTypes[typeObj.Name] = true

		// Add description as a comment if present
		if typeObj.Description != "" {
			lines := strings.Split(typeObj.Description, "\n")
			for _, line := range lines {
				sb.WriteString("# " + line + "\n")
			}
		}

		// Handle different kinds of types
		switch typeObj.Kind {
		case "OBJECT":
			sb.WriteString("type " + typeObj.Name)
			
			// Add interfaces
			if len(typeObj.Interfaces) > 0 {
				sb.WriteString(" implements ")
				for i, interf := range typeObj.Interfaces {
					if i > 0 {
						sb.WriteString(" & ")
					}
					sb.WriteString(interf.Name)
				}
			}
			
			sb.WriteString(" {\n")
			
			// Add fields
			for _, field := range typeObj.Fields {
				// Add field description as comment if present
				if field.Description != "" {
					lines := strings.Split(field.Description, "\n")
					for _, line := range lines {
						sb.WriteString("  # " + line + "\n")
					}
				}
				
				sb.WriteString("  " + field.Name)
				
				// Add arguments if any
				if len(field.Args) > 0 {
					sb.WriteString("(")
					for i, arg := range field.Args {
						if i > 0 {
							sb.WriteString(", ")
						}
						sb.WriteString(arg.Name + ": " + typeRefToString(arg.Type))
						if arg.DefaultValue != "" {
							sb.WriteString(" = " + arg.DefaultValue)
						}
					}
					sb.WriteString(")")
				}
				
				// Add field type
				sb.WriteString(": " + typeRefToString(field.Type) + "\n")
			}
			
			sb.WriteString("}\n\n")
			
		case "INTERFACE":
			sb.WriteString("interface " + typeObj.Name + " {\n")
			
			// Add fields
			for _, field := range typeObj.Fields {
				// Add field description as comment if present
				if field.Description != "" {
					lines := strings.Split(field.Description, "\n")
					for _, line := range lines {
						sb.WriteString("  # " + line + "\n")
					}
				}
				
				sb.WriteString("  " + field.Name)
				
				// Add arguments if any
				if len(field.Args) > 0 {
					sb.WriteString("(")
					for i, arg := range field.Args {
						if i > 0 {
							sb.WriteString(", ")
						}
						sb.WriteString(arg.Name + ": " + typeRefToString(arg.Type))
						if arg.DefaultValue != "" {
							sb.WriteString(" = " + arg.DefaultValue)
						}
					}
					sb.WriteString(")")
				}
				
				// Add field type
				sb.WriteString(": " + typeRefToString(field.Type) + "\n")
			}
			
			sb.WriteString("}\n\n")
			
		case "INPUT_OBJECT":
			sb.WriteString("input " + typeObj.Name + " {\n")
			
			// Add input fields
			for _, field := range typeObj.InputFields {
				// Add field description as comment if present
				if field.Description != "" {
					lines := strings.Split(field.Description, "\n")
					for _, line := range lines {
						sb.WriteString("  # " + line + "\n")
					}
				}
				
				sb.WriteString("  " + field.Name + ": " + typeRefToString(field.Type))
				if field.DefaultValue != "" {
					sb.WriteString(" = " + field.DefaultValue)
				}
				sb.WriteString("\n")
			}
			
			sb.WriteString("}\n\n")
			
		case "ENUM":
			sb.WriteString("enum " + typeObj.Name + " {\n")
			
			// Add enum values
			for _, enumValue := range typeObj.EnumValues {
				// Add value description as comment if present
				if enumValue.Description != "" {
					lines := strings.Split(enumValue.Description, "\n")
					for _, line := range lines {
						sb.WriteString("  # " + line + "\n")
					}
				}
				
				sb.WriteString("  " + enumValue.Name + "\n")
			}
			
			sb.WriteString("}\n\n")
			
		case "UNION":
			sb.WriteString("union " + typeObj.Name + " = ")
			
			// Add possible types
			for i, possibleType := range typeObj.PossibleTypes {
				if i > 0 {
					sb.WriteString(" | ")
				}
				sb.WriteString(possibleType.Name)
			}
			
			sb.WriteString("\n\n")
			
		case "SCALAR":
			sb.WriteString("scalar " + typeObj.Name + "\n\n")
		}
	}

	// Add directives
	for _, directive := range response.Data.Schema.Directives {
		// Skip internal GraphQL directives
		if strings.HasPrefix(directive.Name, "__") {
			continue
		}

		// Add description as a comment if present
		if directive.Description != "" {
			lines := strings.Split(directive.Description, "\n")
			for _, line := range lines {
				sb.WriteString("# " + line + "\n")
			}
		}

		sb.WriteString("directive @" + directive.Name)

		// Add arguments if any
		if len(directive.Args) > 0 {
			sb.WriteString("(")
			for i, arg := range directive.Args {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(arg.Name + ": " + typeRefToString(arg.Type))
				if arg.DefaultValue != "" {
					sb.WriteString(" = " + arg.DefaultValue)
				}
			}
			sb.WriteString(")")
		}

		// Add locations
		sb.WriteString(" on ")
		for i, location := range directive.Locations {
			if i > 0 {
				sb.WriteString(" | ")
			}
			sb.WriteString(location)
		}

		sb.WriteString("\n\n")
	}

	return sb.String()
}

// typeRefToString converts a TypeRef to its string representation in SDL
func typeRefToString(typeRef TypeRef) string {
	if typeRef.Kind == "NON_NULL" && typeRef.OfType != nil {
		return typeRefToString(*typeRef.OfType) + "!"
	} else if typeRef.Kind == "LIST" && typeRef.OfType != nil {
		return "[" + typeRefToString(*typeRef.OfType) + "]"
	} else {
		return typeRef.Name
	}
}
