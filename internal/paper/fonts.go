package paper

import (
	"embed"
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/core/entity"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

//go:embed fonts/*.ttf
var fontsFS embed.FS

// liberationMono is the family name for the embedded Liberation Mono fonts.
const liberationMono = "liberationmono"

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

// loadCustomFonts reads the embedded Liberation Mono fonts and registers all
// styles under the liberationMono font family.
func loadCustomFonts() ([]*entity.CustomFont, error) {
	styles := []struct {
		style fontstyle.Type
		file  string
	}{
		{fontstyle.Normal, "fonts/LiberationMono-Regular.ttf"},
		{fontstyle.Bold, "fonts/LiberationMono-Bold.ttf"},
		{fontstyle.Italic, "fonts/LiberationMono-Italic.ttf"},
		{fontstyle.BoldItalic, "fonts/LiberationMono-BoldItalic.ttf"},
	}

	var fonts []*entity.CustomFont
	for _, s := range styles {
		bytes, err := fontsFS.ReadFile(s.file)
		if err != nil {
			return nil, fmt.Errorf("read embedded font %s: %w", s.file, err)
		}

		fonts = append(fonts, &entity.CustomFont{
			Family: liberationMono,
			Style:  s.style,
			Bytes:  bytes,
		})
	}

	return fonts, nil
}
