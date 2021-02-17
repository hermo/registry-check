package json

import (
	"encoding/json"
	"fmt"

	"github.com/hermo/registry-check/pkg/presentation"
)

// Presenter presents output as JSON
type Presenter struct{}

type jsonPresentation struct {
	ParsedSuccessfully bool     `json:"success"`
	Type               string   `json:"type,omitempty"`
	NumPackages        int      `json:"packages,omitempty"`
	Error              string   `json:"error,omitempty"`
	Filename           string   `json:"filename"`
	Registries         []string `json:"registries"`
}

// NewJSONPresenter creates a new TextPresenter
func NewJSONPresenter() presentation.Presenter {
	return &Presenter{}
}

// Present presents output in text format
func (Presenter) Present(msg *presentation.Output) {
	data := &jsonPresentation{
		ParsedSuccessfully: msg.ParsedSuccessfully,
		Type:               msg.LockfileType,
		NumPackages:        msg.NumPackages,
		Filename:           msg.Filename,
		Registries:         msg.Registries,
	}
	if msg.Error != nil {
		data.Error = msg.Error.Error()
	}
	if b, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		fmt.Println(string(b))
	}
}
