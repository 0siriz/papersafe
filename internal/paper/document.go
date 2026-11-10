package paper

import (
	"fmt"

	"github.com/0siriz/papersafe/internal/encoding/payload"
	"github.com/0siriz/papersafe/internal/encoding/zbase32"
	"github.com/0siriz/papersafe/pkg/keyshard"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

const (
	// documentChunkSize is the maximum number of raw document bytes encoded
	// into a single QR code. Base64 expands this to roughly 1.3 KB, which keeps
	// the generated QR code around version 30 at typical printed sizes.
	documentChunkSize = 1000
	// qrMaxPayloadBytes is the maximum payload size supported by a QR code at
	// error correction level M (version 40). Maroto uses level M internally.
	qrMaxPayloadBytes = 2300
	// documentQRsPerRow is the number of document QR codes placed side by side.
	documentQRsPerRow = 2
)

func GetDocument(quorum keyshard.Quorum, id string, data []byte) (core.Maroto, error) {
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
		text.NewCol(12, segment(zbase32.StdEncoding.EncodeToString(quorum.ID)), props.Text{
			Size:   12,
			Family: liberationMono,
		}),
	)
	m.AddRow(10,
		text.NewCol(12, "Document ID:", props.Text{
			Size:   18,
			Style:  fontstyle.Bold,
			Family: liberationMono,
		}),
	)
	m.AddRow(8,
		text.NewCol(12, id, props.Text{
			Size:   12,
			Family: liberationMono,
		}),
	)

	m.AddRow(10,
		text.NewCol(12,
			"Document data, encrypted with keyshard quorum",
			props.Text{
				Size:   16,
				Color:  &props.Color{Red: 255, Green: 255, Blue: 255},
				Style:  fontstyle.Bold,
				Family: liberationMono,
				Left:   2,
				Top:    2,
			}),
	).WithStyle(&props.Cell{
		BackgroundColor: &props.Color{Red: 204, Green: 104, Blue: 4},
	})

	chunks := chunkDocumentData(data)
	for start := 0; start < len(chunks); start += documentQRsPerRow {
		end := min(start+documentQRsPerRow, len(chunks))

		labels := make([]core.Col, 0, end-start)
		qrs := make([]core.Col, 0, end-start)
		for i := start; i < end; i++ {
			chunk := chunks[i]

			labels = append(labels, text.NewCol(6,
				fmt.Sprintf("Part %d of %d", chunk.Index, chunk.Total),
				props.Text{
					Size:   12,
					Style:  fontstyle.Bold,
					Family: liberationMono,
				}))

			encoded, err := chunk.MarshalText()
			if err != nil {
				return nil, err
			}
			if len(encoded) > qrMaxPayloadBytes {
				return nil, fmt.Errorf(
					"document QR payload %d/%d is too large: %d bytes (max %d)",
					chunk.Index, chunk.Total, len(encoded), qrMaxPayloadBytes,
				)
			}

			qrs = append(qrs, code.NewQrCol(6, string(encoded), props.Rect{Center: true}))
		}

		m.AddRow(6, labels...)
		m.AddRow(80, qrs...)
	}

	return m, nil
}

// chunkDocumentData splits data into indexed payloads carrying at most
// documentChunkSize bytes each. Empty data still produces a single payload so
// the document always has a QR code representing it.
func chunkDocumentData(data []byte) []payload.Payload {
	total := max(1, (len(data)+documentChunkSize-1)/documentChunkSize)

	chunks := make([]payload.Payload, total)
	for i := range total {
		start := i * documentChunkSize
		end := min(start+documentChunkSize, len(data))
		chunks[i] = payload.NewChunk(payload.DomainDocument, i+1, total, data[start:end])
	}

	return chunks
}
