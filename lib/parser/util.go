package parser

import "reflect"

func Coalesce(items ...interface{}) string {
	for _, item := range items {
		if item != nil && reflect.TypeOf(item).Kind() == reflect.String && item != "" {
			return item.(string)
		}
	}

	return ""
}
