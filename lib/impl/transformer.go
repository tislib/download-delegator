package service

import (
	error2 "download-delegator/core/model/errors"
	"download-delegator/lib/parser/model"
	transformers2 "download-delegator/lib/transformers"
	log "github.com/sirupsen/logrus"
)

type TransformerService struct {
	transformers []transformers2.Transformer
}

func (s *TransformerService) Transform(input []byte) ([]byte, error2.State) {
	var data = input
	var err error2.State

	for _, transformer := range s.transformers {
		data, err = transformer.Transform(data)

		if err != error2.NoError {
			return input, err
		}
	}

	return data, error2.NoError
}

func (s *TransformerService) Init(transformerConfigs []model.TransformerConfig) {
	for _, transformConfig := range transformerConfigs {
		transformer := s.locateTransformer(transformConfig.Type)

		if transformer == nil {
			log.Println("could not parse transformer: ", transformConfig.Type)
			continue
		}

		transformer.Init(transformConfig.Parameters)

		s.transformers = append(s.transformers, transformer)
	}
}

func (s *TransformerService) locateTransformer(transformerType model.TransformerType) transformers2.Transformer {
	switch transformerType {
	case model.ScriptTengo:
		return new(transformers2.TengoTransformer)
	case model.ScriptJs:
		return new(transformers2.JsTransformer)
	default:
		return nil
	}
}
