package model

type TransformerType string

const (
	Sanitize    TransformerType = "sanitize"
	Sanitize2   TransformerType = "sanitize2"
	HtmlFormat  TransformerType = "html-format"
	ScriptTengo TransformerType = "script-tengo"
	ScriptJs    TransformerType = "script-js"
)

type TransformerConfig struct {
	Type       TransformerType
	Parameters interface{}
}
