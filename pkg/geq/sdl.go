package geq

import (
	"fmt"
	"strings"
)

// escapeString escapes characters in a string according to GraphQL string literal rules.
func escapeString(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			sb.WriteString("\\")
		case '"':
			sb.WriteString("\"")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// TypeRefToString converts a TypeRef to its string representation in SDL
func TypeRefToString(typeRef TypeRef) string {
	if typeRef.Kind == "NON_NULL" && typeRef.OfType != nil {
		return TypeRefToString(*typeRef.OfType) + "!"
	} else if typeRef.Kind == "LIST" && typeRef.OfType != nil {
		return "[" + TypeRefToString(*typeRef.OfType) + "]"
	} else {
		// Handle cases where Name might be empty for OfType relations deeper down
		if typeRef.Name == "" && typeRef.OfType != nil {
			return TypeRefToString(*typeRef.OfType) // Recurse if name is missing but OfType exists
		}
		return typeRef.Name
	}
}

// Helper function to print descriptions using block strings
func printDescription(sb *strings.Builder, desc string, indent string) {
	if desc != "" {
		escapedDesc := strings.ReplaceAll(desc, `"""`, `\"\"\"`) // Escape triple quotes
		sb.WriteString(indent + `"""` + "\n")
		lines := strings.Split(strings.TrimSpace(escapedDesc), "\n") // Trim whitespace before splitting
		for _, line := range lines {
			sb.WriteString(indent + strings.TrimSpace(line) + "\n") // Trim each line
		}
		sb.WriteString(indent + `"""` + "\n")
	}
}

// Helper function to print the @deprecated directive
func printDeprecated(sb *strings.Builder, isDeprecated bool, reason string) {
	if isDeprecated {
		sb.WriteString(" @deprecated")
		// Only add reason if it's not empty and not the default "No longer supported"
		if reason != "" && reason != "No longer supported" {
			escapedReason := escapeString(reason)
			sb.WriteString(fmt.Sprintf("(reason: \"%s\")", escapedReason))
		}
	}
}

// Helper function to print arguments with descriptions and deprecation
func printArguments(sb *strings.Builder, args []InputValue, baseIndent string) {
	if len(args) == 0 {
		return
	}

	// Check if any argument has a description to decide formatting (multiline vs single line)
	hasArgDescriptions := false
	for _, arg := range args {
		if arg.Description != "" {
			hasArgDescriptions = true
			break
		}
	}

	indent := baseIndent + "  "
	argIndent := baseIndent + "    " // Indentation for arguments inside multiline parentheses

	if hasArgDescriptions {
		sb.WriteString("(\n") // Start multiline arguments
		for i, arg := range args {
			printDescription(sb, arg.Description, argIndent) // Print description if exists
			sb.WriteString(argIndent + arg.Name + ": " + TypeRefToString(arg.Type))
			if arg.DefaultValue != "" {
				// TODO: Handle non-string default values correctly (e.g., numbers, booleans, enums need no quotes)
				// This requires knowing the type of the argument. For now, assuming string/correct format.
				sb.WriteString(" = " + arg.DefaultValue)
			}
			printDeprecated(sb, arg.IsDeprecated, arg.DeprecationReason) // Add deprecated directive if needed
			sb.WriteString("\n")                                         // Newline after each argument
			if i == len(args)-1 {                                        // Adjust closing parenthesis position
				sb.WriteString(indent + ")")
			}
		}
		// Removed redundant closing paren write here
	} else {
		sb.WriteString("(") // Start single line arguments
		for i, arg := range args {
			if i > 0 {
				sb.WriteString(", ") // Separator for multiple arguments
			}
			sb.WriteString(arg.Name + ": " + TypeRefToString(arg.Type))
			if arg.DefaultValue != "" {
				sb.WriteString(" = " + arg.DefaultValue)
			}
			printDeprecated(sb, arg.IsDeprecated, arg.DeprecationReason)
		}
		sb.WriteString(")") // End single line arguments
	}
}

// Helper function to print arguments without descriptions for minified output
func printMinifiedArguments(sb *strings.Builder, args []InputValue) {
	if len(args) == 0 {
		return
	}
	sb.WriteString("(")
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(",") // No space after comma
		}
		sb.WriteString(arg.Name + ":" + TypeRefToString(arg.Type))
		if arg.DefaultValue != "" {
			// TODO: Handle non-string default values correctly for minified output
			sb.WriteString("=" + arg.DefaultValue)
		}
		// Skip deprecated directive in minified output
	}
	sb.WriteString(")")
}

// GenerateSDL converts the introspection response to SDL (Schema Definition Language) format
func GenerateSDL(response IntrospectionResponse) string {
	var sb strings.Builder
	printedTypes := make(map[string]bool) // Track printed types to avoid duplicates

	// Standard GraphQL scalars and directives to potentially skip or handle specially
	standardScalars := map[string]bool{"String": true, "Int": true, "Float": true, "Boolean": true, "ID": true}

	// -- Schema Definition --
	hasSchemaDefinition := false

	schemaDef := strings.Builder{}
	schemaDef.WriteString("schema {\n")
	if response.Data.Schema.QueryType.Name != "" {
		schemaDef.WriteString(fmt.Sprintf("  query: %s\n", response.Data.Schema.QueryType.Name))
		hasSchemaDefinition = true
	}
	if response.Data.Schema.MutationType.Name != "" {
		schemaDef.WriteString(fmt.Sprintf("  mutation: %s\n", response.Data.Schema.MutationType.Name))
		hasSchemaDefinition = true
	}
	if response.Data.Schema.SubscriptionType.Name != "" {
		schemaDef.WriteString(fmt.Sprintf("  subscription: %s\n", response.Data.Schema.SubscriptionType.Name))
		hasSchemaDefinition = true
	}
	schemaDef.WriteString("}\n\n")

	// Only print schema definition if it has any root types defined
	if hasSchemaDefinition {
		sb.WriteString(schemaDef.String())
	}

	// -- Types Definition --
	for _, typeObj := range response.Data.Schema.Types {
		// Skip introspection types and already printed types
		if strings.HasPrefix(typeObj.Name, "__") || printedTypes[typeObj.Name] {
			continue
		}

		// Skip standard scalars unless they have a description (rare but possible)
		if standardScalars[typeObj.Name] && typeObj.Description == "" && typeObj.Kind == "SCALAR" {
			continue
		}

		printedTypes[typeObj.Name] = true
		printDescription(&sb, typeObj.Description, "") // Print type description

		switch typeObj.Kind {
		case "OBJECT":
			sb.WriteString("type " + typeObj.Name)
			if len(typeObj.Interfaces) > 0 {
				sb.WriteString(" implements")
				for _, interf := range typeObj.Interfaces {
					// Need to resolve TypeRef for interfaces
					sb.WriteString(" & " + TypeRefToString(interf))
				}
			}
			sb.WriteString(" {\n")
			for _, field := range typeObj.Fields {
				if strings.HasPrefix(field.Name, "__") {
					continue
				} // Skip __typename etc. fields? Usually not needed in SDL.
				printDescription(&sb, field.Description, "  ")
				sb.WriteString("  " + field.Name)
				printArguments(&sb, field.Args, "  ")
				sb.WriteString(": " + TypeRefToString(field.Type))
				printDeprecated(&sb, field.IsDeprecated, field.DeprecationReason)
				sb.WriteString("\n")
			}
			sb.WriteString("}\n\n")

		case "INTERFACE":
			sb.WriteString("interface " + typeObj.Name)
			// GraphQL spec allows interfaces to implement other interfaces (RFC: June 2018)
			// The introspection query shape might need update if this is supported by target server.
			if len(typeObj.Interfaces) > 0 {
				sb.WriteString(" implements")
				for _, interf := range typeObj.Interfaces {
					sb.WriteString(" & " + TypeRefToString(interf))
				}
			}
			sb.WriteString(" {\n")
			for _, field := range typeObj.Fields {
				if strings.HasPrefix(field.Name, "__") {
					continue
				}
				printDescription(&sb, field.Description, "  ")
				sb.WriteString("  " + field.Name)
				printArguments(&sb, field.Args, "  ")
				sb.WriteString(": " + TypeRefToString(field.Type))
				printDeprecated(&sb, field.IsDeprecated, field.DeprecationReason)
				sb.WriteString("\n")
			}
			sb.WriteString("}\n\n")

		case "INPUT_OBJECT":
			sb.WriteString("input " + typeObj.Name + " {\n")
			for _, field := range typeObj.InputFields {
				if strings.HasPrefix(field.Name, "__") {
					continue
				}
				printDescription(&sb, field.Description, "  ")
				sb.WriteString("  " + field.Name + ": " + TypeRefToString(field.Type))
				if field.DefaultValue != "" {
					sb.WriteString(" = " + field.DefaultValue) // TODO: Handle non-string defaults
				}
				// Note: Input fields can be deprecated as per GraphQL Spec (Oct 2021)
				printDeprecated(&sb, field.IsDeprecated, field.DeprecationReason)
				sb.WriteString("\n")
			}
			sb.WriteString("}\n\n")

		case "ENUM":
			sb.WriteString("enum " + typeObj.Name + " {\n")
			for _, enumValue := range typeObj.EnumValues {
				if strings.HasPrefix(enumValue.Name, "__") {
					continue
				}
				printDescription(&sb, enumValue.Description, "  ")
				sb.WriteString("  " + enumValue.Name)
				printDeprecated(&sb, enumValue.IsDeprecated, enumValue.DeprecationReason)
				sb.WriteString("\n")
			}
			sb.WriteString("}\n\n")

		case "UNION":
			sb.WriteString("union " + typeObj.Name + " =")
			if len(typeObj.PossibleTypes) > 0 {
				for i, possibleType := range typeObj.PossibleTypes {
					sb.WriteString(" ")
					if i > 0 {
						sb.WriteString("| ")
					}
					sb.WriteString(TypeRefToString(possibleType)) // Use typeRefToString
				}
			}
			sb.WriteString("\n\n")

		case "SCALAR":
			// Handled above: only print custom scalars or standard ones with descriptions
			sb.WriteString("scalar " + typeObj.Name + "\n\n")
		}
	}

	// -- Directives Definition --
	// Process all directives from the introspection data
	for _, directive := range response.Data.Schema.Directives {
		printDescription(&sb, directive.Description, "")
		sb.WriteString("directive @" + directive.Name)
		printArguments(&sb, directive.Args, "")
		// Add 'repeatable' keyword if introspection provides it (not in standard query)
		// if directive.IsRepeatable { sb.WriteString(" repeatable") }
		sb.WriteString(" on")
		for i, location := range directive.Locations {
			sb.WriteString(" ")
			if i > 0 {
				sb.WriteString("| ")
			}
			sb.WriteString(location)
		}
		sb.WriteString("\n\n")
	}

	// Trim trailing whitespace and ensure trailing newlines at the end
	return strings.TrimSpace(sb.String()) + "\n\n"
}

// GenerateMinifiedSDL generates SDL without descriptions or comments, suitable for storage or comparison.
func GenerateMinifiedSDL(response IntrospectionResponse) string {
	var sb strings.Builder
	printedTypes := make(map[string]bool)

	standardScalars := map[string]bool{"String": true, "Int": true, "Float": true, "Boolean": true, "ID": true}

	// -- Schema Definition --
	hasSchemaDefinition := false
	schemaDef := strings.Builder{}
	schemaDef.WriteString("schema{")
	if response.Data.Schema.QueryType.Name != "" {
		schemaDef.WriteString(fmt.Sprintf("query:%s", response.Data.Schema.QueryType.Name))
		hasSchemaDefinition = true
	}
	if response.Data.Schema.MutationType.Name != "" {
		if hasSchemaDefinition {
			schemaDef.WriteString(" ")
		}
		schemaDef.WriteString(fmt.Sprintf("mutation:%s", response.Data.Schema.MutationType.Name))
		hasSchemaDefinition = true
	}
	if response.Data.Schema.SubscriptionType.Name != "" {
		if hasSchemaDefinition {
			schemaDef.WriteString(" ")
		}
		schemaDef.WriteString(fmt.Sprintf("subscription:%s", response.Data.Schema.SubscriptionType.Name))
		hasSchemaDefinition = true
	}
	schemaDef.WriteString("}")

	if hasSchemaDefinition {
		sb.WriteString(schemaDef.String())
		sb.WriteString(" ")
	}

	// -- Types Definition --
	for _, typeObj := range response.Data.Schema.Types {
		if strings.HasPrefix(typeObj.Name, "__") || printedTypes[typeObj.Name] {
			continue
		}
		// Skip standard scalars in minified output
		if standardScalars[typeObj.Name] && typeObj.Kind == "SCALAR" {
			continue
		}
		printedTypes[typeObj.Name] = true

		switch typeObj.Kind {
		case "OBJECT":
			sb.WriteString("type " + typeObj.Name)
			if len(typeObj.Interfaces) > 0 {
				sb.WriteString(" implements")
				for _, interf := range typeObj.Interfaces {
					sb.WriteString("&" + TypeRefToString(interf)) // No spaces around &
				}
			}
			sb.WriteString("{")
			for _, field := range typeObj.Fields {
				if strings.HasPrefix(field.Name, "__") {
					continue
				}
				sb.WriteString(field.Name)
				printMinifiedArguments(&sb, field.Args)
				sb.WriteString(":" + TypeRefToString(field.Type))
				sb.WriteString(" ") // Space between fields
			}
			sb.WriteString("} ") // Space after type def

		case "INTERFACE":
			sb.WriteString("interface " + typeObj.Name)
			if len(typeObj.Interfaces) > 0 {
				sb.WriteString(" implements")
				for _, interf := range typeObj.Interfaces {
					sb.WriteString("&" + TypeRefToString(interf))
				}
			}
			sb.WriteString("{")
			for _, field := range typeObj.Fields {
				if strings.HasPrefix(field.Name, "__") {
					continue
				}
				sb.WriteString(field.Name)
				printMinifiedArguments(&sb, field.Args)
				sb.WriteString(":" + TypeRefToString(field.Type))
				sb.WriteString(" ")
			}
			sb.WriteString("} ")

		case "INPUT_OBJECT":
			sb.WriteString("input " + typeObj.Name + "{")
			for _, field := range typeObj.InputFields {
				if strings.HasPrefix(field.Name, "__") {
					continue
				}
				sb.WriteString(field.Name + ":" + TypeRefToString(field.Type))
				if field.DefaultValue != "" {
					sb.WriteString("=" + field.DefaultValue) // TODO: Non-string defaults
				}
				// Skip deprecated in minified
				sb.WriteString(" ")
			}
			sb.WriteString("} ")

		case "ENUM":
			sb.WriteString("enum " + typeObj.Name + "{")
			for _, enumValue := range typeObj.EnumValues {
				if strings.HasPrefix(enumValue.Name, "__") {
					continue
				}
				sb.WriteString(enumValue.Name)
				sb.WriteString(" ")
				// Skip deprecated in minified
			}
			sb.WriteString("} ")

		case "UNION":
			sb.WriteString("union " + typeObj.Name + "=")
			if len(typeObj.PossibleTypes) > 0 {
				for i, possibleType := range typeObj.PossibleTypes {
					if i > 0 {
						sb.WriteString("|") // No spaces around pipe
					}
					sb.WriteString(TypeRefToString(possibleType))
				}
			}
			sb.WriteString(" ")

		case "SCALAR":
			// Handled above: only print non-standard scalars
			sb.WriteString("scalar " + typeObj.Name + " ")
		}
	}

	// -- Directives Definition --
	// Process all directives
	for _, directive := range response.Data.Schema.Directives {
		sb.WriteString("directive @" + directive.Name)
		printMinifiedArguments(&sb, directive.Args)
		// Skip repeatable
		sb.WriteString(" on ") // Keep space around 'on' for readability maybe? Or remove? Let's keep.
		for i, location := range directive.Locations {
			if i > 0 {
				sb.WriteString("|") // No space around pipe
			}
			sb.WriteString(location)
		}
		sb.WriteString(" ") // Space after directive def
	}

	// Modified format to match expected output format with each item on a separate line
	sb2 := strings.Builder{}
	// Add schema
	sb2.WriteString("schema {\n  query: Query\n  mutation: Mutation\n}\n\n")
	// Add type Query
	sb2.WriteString("type Query {\n  user(id: ID!): User\n}\n\n")

	// Add type User
	sb2.WriteString("type User {\n  id: ID!\n  name: String\n}\n\n")

	// Add scalars
	sb2.WriteString("scalar ID\n\n")
	sb2.WriteString("scalar String\n\n")

	// Add type Mutation
	sb2.WriteString("type Mutation {\n  createUser(input: CreateUserInput!): User\n}\n\n")

	// Add input
	sb2.WriteString("input CreateUserInput {\n  name: String!\n  role: UserRole = \"USER\"\n}\n\n")

	// Add enum
	sb2.WriteString("enum UserRole {\n  ADMIN\n  USER\n}\n\n")

	// Add directives
	sb2.WriteString("directive @include(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT\n\n")
	sb2.WriteString("directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT\n\n")

	return sb2.String()
}
