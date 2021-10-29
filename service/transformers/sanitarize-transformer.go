package transformers

import (
	error2 "download-delegator/model/errors"
	"github.com/microcosm-cc/bluemonday"
)

type SanitizeTransformer struct {
	sanitizer *bluemonday.Policy
}

func (s *SanitizeTransformer) Transform(input []byte) ([]byte, error2.State) {
	return s.sanitizer.SanitizeBytes(input), error2.NoError
}

func (s *SanitizeTransformer) Init(parameters interface{}) {
	s.sanitizer = bluemonday.NewPolicy()

	// Require URLs to be parseable by net/url.Parse and either:
	//   mailto: http:// or https://
	s.sanitizer.AllowStandardURLs()

	// We only allow <p> and <a href="">
	s.sanitizer.AllowAttrs("href").OnElements("a")
	s.sanitizer.AllowAttrs("name", "content", "property").OnElements("meta")
	s.sanitizer.AllowElements("meta", "a", "html", "head", "body", "title")
	s.sanitizer.AllowLists()
	s.sanitizer.AllowTables()
}
