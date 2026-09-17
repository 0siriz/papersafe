package paper

import (
	"strconv"

	"github.com/0siriz/papersafe/pkg/keyshard"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func GetKeyshard(shard keyshard.KeyShard) (core.Maroto, error) {
	m, err := newMaroto()
	if err != nil {
		return nil, err
	}

	m.AddRow(10,
		text.NewCol(12, "Quorum ID:", props.Text{
			Size:   18,
			Style:  fontstyle.Bold,
			Family: liberationMono,
		}),
	)
	m.AddRow(8,
		text.NewCol(12, segment(zbase32.EncodeToString(shard.PublicKey)), props.Text{
			Size:   12,
			Family: liberationMono,
		}),
	)
	m.AddRow(10,
		text.NewCol(12, "Shard ID:", props.Text{
			Size:   18,
			Style:  fontstyle.Bold,
			Family: liberationMono,
		}),
	)
	m.AddRow(8,
		text.NewCol(12, strconv.Itoa(int(shard.ID&0xFF)), props.Text{
			Size:   12,
			Family: liberationMono,
		}),
	)

	m.AddRow(10,
		text.NewCol(12,
			"Key shard data, encrypted with mnemonic",
			props.Text{
				Size:   16,
				Color:  &props.Color{Red: 255, Green: 255, Blue: 255},
				Style:  fontstyle.Bold,
				Family: liberationMono,
				Left:   2,
				Top:    2,
			}),
	).WithStyle(&props.Cell{
		BackgroundColor: &props.Color{Red: 9, Green: 158, Blue: 86},
	})

	shardData, err := shard.MarshalBinary()
	if err != nil {
		return nil, err
	}

	m.AddRow(80,
		code.NewQrCol(6,
			string(shardData),
			props.Rect{
				Center: true,
			}),
		text.NewCol(6, segment(zbase32.EncodeToString(shardData)), props.Text{
			Size:   12,
			Family: liberationMono,
		}),
	)
	m.AddRow(10,
		text.NewCol(12,
			"Mnemonic",
			props.Text{
				Size:   16,
				Color:  &props.Color{Red: 255, Green: 255, Blue: 255},
				Style:  fontstyle.Bold,
				Family: liberationMono,
				Left:   2,
				Top:    2,
			}),
	).WithStyle(&props.Cell{
		BackgroundColor: &props.Color{Red: 9, Green: 158, Blue: 86},
	})

	for range 8 {
		m.AddRow(15,
			func() []core.Col {
				cols := make([]core.Col, 3)
				for i := range 3 {
					cols[i] = line.NewCol(4, props.Line{OffsetPercent: 100})
				}
				return cols
			}()...,
		)
	}

	return m, nil
}
