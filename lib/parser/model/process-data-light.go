package model

type ProcessDataLight struct {
	Html *string `json:"html" binding:"required"`
	Url  *string `json:"url" binding:"required"`
}
