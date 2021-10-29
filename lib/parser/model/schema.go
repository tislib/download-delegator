package model

type PropertiesType map[string]SchemaProperty

type Schema struct {
	Version     string         `json:"version" binding:"required"`
	Namespace   string         `json:"namespace"`
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Tags        []string       `json:"tags"`
	Properties  PropertiesType `json:"properties" binding:"required"`
}

type SchemaProperty struct {
	// common fields
	Description  string `json:"description" binding:"required"`
	IsAllowQuery bool   `json:"isAllowQuery" binding:"required"`
	Example      string `json:"example" binding:"required"`
	// type
	Type string `json:"type" binding:"required"`
	// string
	Pattern string `json:"pattern" binding:"required"`
	// array
	Items *SchemaProperty `json:"items" binding:"required"`
	// object
	Properties PropertiesType `json:"properties" binding:"required"`
	// reference
	Schema string `json:"schema" binding:"required"`
}

func (schemaProperty SchemaProperty) IsArrayProperty() bool {
	return schemaProperty.Type == "array"
}

func (schemaProperty SchemaProperty) IsObjectProperty() bool {
	return schemaProperty.Type == "object"
}

func (schemaProperty SchemaProperty) IsReferenceProperty() bool {
	return schemaProperty.Type == "ref"
}

func (schemaProperty SchemaProperty) IsStringProperty() bool {
	return schemaProperty.Type == "string"
}

func (schemaProperty SchemaProperty) IsNumberProperty() bool {
	return schemaProperty.Type == "number"
}

type ObjectProperty interface {
	GetProperties() PropertiesType
}

func (schema *Schema) GetProperties() PropertiesType {
	return schema.Properties
}

func (schemaProperty *SchemaProperty) GetProperties() PropertiesType {
	return schemaProperty.Properties
}
