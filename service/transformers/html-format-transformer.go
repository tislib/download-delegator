package transformers

import (
	error2 "download-delegator/model/errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yosssi/gohtml"
)

type HtmlFormatTransformer struct {
	sanitizer *bluemonday.Policy
}

func (s *HtmlFormatTransformer) Transform(input []byte) ([]byte, error2.State) {
	return gohtml.FormatBytes(input), error2.NoError
}

func (s *HtmlFormatTransformer) Init(parameters interface{}) {

}
