package transformers

import (
	"download-delegator/lib/parser"
	"download-delegator/lib/parser/model"
	error2 "download-delegator/model/errors"
	"github.com/dop251/goja"
	"log"
)

type JsTransformer struct {
	prog *goja.Program
	lib  struct {
		ParseHtml func(htmlContent string) model.DocNode
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
	s.lib.ParseHtml = func(htmlContent string) model.DocNode {
		res, err := parser.ParseHtml(htmlContent)

		if err != nil {
			log.Print(err)

			return nil
		}

		return res
	}

	prog, err := goja.Compile("js-transformer", parameters.(string), false)

	if err != nil {
		log.Print(err)
	}

	s.prog = prog
}
