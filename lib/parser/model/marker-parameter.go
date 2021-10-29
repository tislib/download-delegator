package model

type MarkerParameter struct {
	Name          string                 `json:"name"`
	Caption       string                 `json:"caption"`
	ParameterType ParameterType          `json:"parameterType"`
	Options       MarkerParameterOptions `json:"options"`
	Values        []interface{}          `json:"values"`
	DefaultValue  interface{}            `json:"defaultValue"`
	Required      bool                   `json:"required"`
}

type ParameterType string

type MarkerParameterOptions map[string]interface{}

const (
	TEXT               ParameterType = "TEXT"
	REGEXP                           = "REGEXP"
	REGEX_SUB                        = "REGEX_SUB"
	INSPECTOR                        = "INSPECTOR"
	NUMBER                           = "NUMBER"
	COMBOBOX                         = "COMBOBOX"
	CHECKBOX                         = "CHECKBOX"
	TRANSITIVE_ELEMENT ParameterType = "TEXT"
)
