package model

type TransformerType string

const (
	ScriptTengo TransformerType = "script-tengo"
	ScriptJs    TransformerType = "script-js"
)

type TransformerConfig struct {
	Type       TransformerType
	Parameters interface{}
}
