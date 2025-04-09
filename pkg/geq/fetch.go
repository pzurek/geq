package geq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchIntrospectionJSON fetches the GraphQL schema using the standard introspection query.
// It takes the GraphQL endpoint URL and an optional header string (e.g., "Authorization: Bearer token").
// It returns the raw JSON response as a string.
func FetchIntrospectionJSON(endpoint, headerStr string) (string, error) {
	// Use the canonical introspection query from graphql-js
	introspectionQuery := `
    query IntrospectionQuery {
      __schema {
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
          ...FullType
        }
        directives {
          name
          description
          locations
          args {
            ...InputValue
          }
        }
      }
    }

    fragment FullType on __Type {
      kind
      name
      description
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }

    fragment InputValue on __InputValue {
      name
      description
      type { ...TypeRef }
      defaultValue
    }

    fragment TypeRef on __Type {
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
  `

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
		// Try to unmarshal the error response for better formatting, fallback to raw body
		var errorResp struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}
		if json.Unmarshal(body, &errorResp) == nil && len(errorResp.Errors) > 0 {
			errorMessages := make([]string, len(errorResp.Errors))
			for i, e := range errorResp.Errors {
				errorMessages[i] = e.Message
			}
			return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, strings.Join(errorMessages, "; "))
		}
		return "", fmt.Errorf("server returned error: %s", body) // Fallback to raw body
	}

	return string(body), nil // Return raw JSON string
}
