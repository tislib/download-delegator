package parser

import (
	"github.com/pborman/uuid"
)

func NamedUUID(data []byte) string {
	return uuid.NewMD5(nil, data).String()
}
