package model

type ProcessData struct {
	Model            *Model  `json:"model" binding:"required"`
	AdditionalModels []Model `json:"additionalModels"`
	Schema           *Schema `json:"schema" binding:"required"`

	Html *string `json:"html" binding:"required"`
	Url  *string `json:"url" binding:"required"`
}
