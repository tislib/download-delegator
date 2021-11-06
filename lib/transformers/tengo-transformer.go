package transformers

import (
	"context"
	error2 "download-delegator/core/model/errors"
	"github.com/d5/tengo/v2"
	log "github.com/sirupsen/logrus"
)

type TengoTransformer struct {
	c *tengo.Compiled
}

func (s *TengoTransformer) Transform(input []byte) ([]byte, error2.State) {
	if s.c == nil {
		return input, error2.TengoCompileError
	}

	compiled := s.c.Clone()

	compiled.Set("input", input)

	err := compiled.RunContext(context.Background())

	if err != nil {
		log.Println(err)
	}

	return compiled.Get("output").Bytes(), error2.NoError
}

func (s *TengoTransformer) Init(parameters interface{}) {
	script := tengo.NewScript([]byte(parameters.(string)))

	var func1 = func(args ...tengo.Object) (ret tengo.Object, err error) {
		return nil, nil
	}

	err := script.Add("func1", func1)
	script.Add("input", []byte{})

	if err != nil {
		log.Print(err)
	}

	s.c, err = script.Compile()

	if err != nil {
		log.Print(err)
	}
}
