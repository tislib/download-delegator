package transformers

import (
	error2 "download-delegator/core/model/errors"
	"download-delegator/lib/parser"
	"download-delegator/lib/parser/model"
	"github.com/dop251/goja"
	log "github.com/sirupsen/logrus"
)

type JsTransformer struct {
	prog *goja.Program
	lib  struct {
		Parser struct {
			ParseHtml func(htmlContent string) (model.DocNode, error)
			Processor parser.Processor
		}
	}
}

func (s *JsTransformer) Transform(input []byte) ([]byte, error2.State) {
	if s.prog == nil {
		return input, error2.JsCompileError
	}

	runtime := goja.New()

	err := runtime.Set("input", string(input))

	if err != nil {
		log.Print(err)
	}

	err = runtime.Set("lib", s.lib)

	if err != nil {
		log.Print(err)
	}

	result, err := runtime.RunProgram(s.prog)

	if err != nil {
		log.Print(err)

		return input, error2.JsCompileError
	}

	return []byte(result.Export().(string)), error2.NoError
}

func (s *JsTransformer) Init(parameters interface{}) {
	s.initLib()

	prog, err := goja.Compile("js-transformer", parameters.(string), false)

	if err != nil {
		log.Print(err)
	}

	s.prog = prog
}

func (s *JsTransformer) initLib() {
	s.lib.Parser.ParseHtml = parser.ParseHtml
	s.lib.Parser.Processor = parser.Processor{}
}
