package transformers

import (
	error2 "download-delegator/model/errors"
	"github.com/microcosm-cc/bluemonday"
)

type Sanitize2Transformer struct {
	sanitizer *bluemonday.Policy
}

func (s *Sanitize2Transformer) Transform(input []byte) ([]byte, error2.State) {
	return s.sanitizer.SanitizeBytes(input), error2.NoError
}

func (s *Sanitize2Transformer) Init(parameters interface{}) {
	s.sanitizer = bluemonday.NewPolicy()

	s.sanitizer.AllowAttrs("name", "content", "property").OnElements("meta")
	s.sanitizer.AllowElements("meta", "html", "head", "title")
	s.sanitizer.SkipElementsContent("body")
}
