package geq

// InputValue represents a GraphQL input value definition
type InputValue struct {
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Type              TypeRef `json:"type"`
	DefaultValue      string  `json:"defaultValue"`
	IsDeprecated      bool    `json:"isDeprecated"`
	DeprecationReason string  `json:"deprecationReason"`
}

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
					Name              string       `json:"name"`
					Description       string       `json:"description"`
					Args              []InputValue `json:"args"`
					Type              TypeRef      `json:"type"`
					IsDeprecated      bool         `json:"isDeprecated"`
					DeprecationReason string       `json:"deprecationReason"`
				} `json:"fields"`
				InputFields []struct {
					Name              string  `json:"name"`
					Description       string  `json:"description"`
					Type              TypeRef `json:"type"`
					DefaultValue      string  `json:"defaultValue"`
					IsDeprecated      bool    `json:"isDeprecated"`
					DeprecationReason string  `json:"deprecationReason"`
				} `json:"inputFields"`
				Interfaces []TypeRef `json:"interfaces"`
				EnumValues []struct {
					Name              string `json:"name"`
					Description       string `json:"description"`
					IsDeprecated      bool   `json:"isDeprecated"`
					DeprecationReason string `json:"deprecationReason"`
				} `json:"enumValues"`
				PossibleTypes []TypeRef `json:"possibleTypes"`
			} `json:"types"`
			Directives []struct {
				Name        string       `json:"name"`
				Description string       `json:"description"`
				Locations   []string     `json:"locations"`
				Args        []InputValue `json:"args"`
			} `json:"directives"`
		} `json:"__schema"`
	} `json:"data"`
}

// TypeRef represents a GraphQL type reference
type TypeRef struct {
	Kind   string   `json:"kind"`
	Name   string   `json:"name"`
	OfType *TypeRef `json:"ofType"`
}
