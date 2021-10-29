package transformers

import error2 "download-delegator/model/errors"

type Transformer interface {
	Init(parameters interface{})
	Transform(input []byte) ([]byte, error2.State)
}
