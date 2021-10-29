package model

import (
	"reflect"
	"strings"
)

type Marker interface {
	Apply(data *PageData, markerData MarkerData) *PageData
	GetName() string
	GetParameters() []MarkerParameter
}

func GetMarkerStringParameter(markerData MarkerData, paramName string) string {
	val := markerData.Parameters[paramName]

	if reflect.TypeOf(val).Kind() == reflect.String {
		return val.(string)
	}

	return ""
}

func GetMarkerBooleanParameter(markerData MarkerData, paramName string) bool {
	val := markerData.Parameters[paramName]

	if val == nil {
		return false
	}

	if reflect.TypeOf(val).Kind() == reflect.String {
		return strings.ToLower(val.(string)) == "true"
	} else if reflect.TypeOf(val).Kind() == reflect.Bool {
		return val.(bool)
	}

	return false
}
