package paper

import (
	"strings"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

const (
	segmentSize    = 5
	segmentSpacing = 1
)

func segment(text string) string {
	var result strings.Builder

	for i := 0; i < len(text); i += segmentSize {
		end := min(i+segmentSize, len(text))

		result.WriteString(text[i:end])
		result.WriteByte(0x20)
	}

	return result.String()
}

// newMaroto creates a Maroto instance using the embedded Liberation Mono fonts
// as the default font family.
func newMaroto() (core.Maroto, error) {
	customFonts, err := loadCustomFonts()
	if err != nil {
		return nil, err
	}

	cfg := config.NewBuilder().
		WithCustomFonts(customFonts).
		WithDefaultFont(&props.Font{Family: liberationMono}).
		Build()

	return maroto.New(cfg), nil
}
