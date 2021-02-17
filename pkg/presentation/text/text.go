package text

import (
	"fmt"
	"os"

	"github.com/hermo/registry-check/pkg/presentation"
)

// Presenter presents output as text
type Presenter struct{}

// NewTextPresenter creates a new TextPresenter
func NewTextPresenter() presentation.Presenter {
	return &Presenter{}
}

// Present presents output in text format
func (Presenter) Present(msg *presentation.Output) {
	if msg.ParsedSuccessfully {
		fmt.Printf("Found %d registries in %s with %d packages:\n", len(msg.Registries), msg.Filename, msg.NumPackages)
		for _, reg := range msg.Registries {
			fmt.Printf("- %s\n", reg)
		}
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", msg.Filename, msg.Error)
	}
}
