package service

import (
	"download-delegator/lib/parser/model"
	error2 "download-delegator/model/errors"
	"download-delegator/service/transformers"
	"log"
)

type TransformerService struct {
	transformers []transformers.Transformer
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

func (s *TransformerService) locateTransformer(transformerType model.TransformerType) transformers.Transformer {
	switch transformerType {
	case model.Sanitize:
		return new(transformers.SanitizeTransformer)
	case model.Sanitize2:
		return new(transformers.Sanitize2Transformer)
	case model.HtmlFormat:
		return new(transformers.HtmlFormatTransformer)
	case model.ScriptTengo:
		return new(transformers.TengoTransformer)
	case model.ScriptJs:
		return new(transformers.JsTransformer)
	default:
		return nil
	}
}
