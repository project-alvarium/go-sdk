package metadata

import "github.com/project-alvarium/go-sdk/pkg/annotator"

const kind = "iota"

type Instance struct {
	kind             string
	Result           string   `json:"result"`
	ValidTransaction bool     `json:"validTransaction"`
	Unique           []string `json:"unique"`
}

func New(validTransaction bool, unique []string) *Instance {
	return &Instance{
		kind:             kind,
		Result:           annotator.SuccessKind,
		ValidTransaction: validTransaction,
		Unique:           unique,
	}
}

func (*Instance) Kind() string {
	return kind
}