package model

import "github.com/google/uuid"

type RecordData map[string]Value
type RecordMeta map[string]string

type Value interface{}

type Record struct {
	Id          uuid.UUID  `json:"id" binding:"required"`
	Tags        []string   `json:"tags"`
	Description string     `json:"description"`
	Data        RecordData `json:"data" binding:"required"`
	Meta        RecordMeta `json:"meta"`
	Schema      string     `json:"schema" binding:"required"`

	// reference properties
	Ref        string `json:"ref" binding:"required"`
	Source     string `json:"source" binding:"required"`
	SourceUrl  string `json:"sourceUrl" binding:"required"`
	ObjectType string `json:"objectType" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type Reference struct {
	Ref        string `json:"ref" binding:"required"`
	Source     string `json:"source" binding:"required"`
	SourceUrl  string `json:"sourceUrl" binding:"required"`
	ObjectType string `json:"objectType" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

//String ref;
//    String source;
//    String sourceUrl;
//    String objectType;
//    String name;
