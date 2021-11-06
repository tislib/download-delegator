package transformers

import (
	error2 "download-delegator/core/model/errors"
)

type Transformer interface {
	Init(parameters interface{})
	Transform(input []byte) ([]byte, error2.State)
}
