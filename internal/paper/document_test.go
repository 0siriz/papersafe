package paper

import (
	"bytes"
	"testing"

	"github.com/0siriz/papersafe/internal/encoding/payload"
)

func TestChunkDocumentData(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		chunks := chunkDocumentData(nil)

		if len(chunks) != 1 {
			t.Fatalf("expected 1 chunk, got %d", len(chunks))
		}
		if chunks[0].Index != 1 || chunks[0].Total != 1 {
			t.Fatalf("unexpected chunk metadata: %d/%d", chunks[0].Index, chunks[0].Total)
		}
		if len(chunks[0].Data) != 0 {
			t.Fatalf("expected empty data, got %d bytes", len(chunks[0].Data))
		}
	})

	t.Run("single chunk", func(t *testing.T) {
		data := bytes.Repeat([]byte{0xab}, documentChunkSize)
		chunks := chunkDocumentData(data)

		if len(chunks) != 1 {
			t.Fatalf("expected 1 chunk, got %d", len(chunks))
		}
		if chunks[0].Index != 1 || chunks[0].Total != 1 {
			t.Fatalf("unexpected chunk metadata: %d/%d", chunks[0].Index, chunks[0].Total)
		}
		if !bytes.Equal(chunks[0].Data, data) {
			t.Fatal("chunk data mismatch")
		}
	})

	t.Run("multiple chunks", func(t *testing.T) {
		data := make([]byte, documentChunkSize*2+123)
		for i := range data {
			data[i] = byte(i)
		}

		chunks := chunkDocumentData(data)
		if len(chunks) != 3 {
			t.Fatalf("expected 3 chunks, got %d", len(chunks))
		}

		var joined []byte
		for i, chunk := range chunks {
			if chunk.Index != i+1 {
				t.Fatalf("chunk %d has index %d", i, chunk.Index)
			}
			if chunk.Total != len(chunks) {
				t.Fatalf("chunk %d has total %d", i, chunk.Total)
			}
			if len(chunk.Data) > documentChunkSize {
				t.Fatalf("chunk %d is too large: %d bytes", i, len(chunk.Data))
			}
			joined = append(joined, chunk.Data...)
		}

		if !bytes.Equal(joined, data) {
			t.Fatal("reassembled data mismatch")
		}
	})

	t.Run("encoded payloads fit qr capacity", func(t *testing.T) {
		data := bytes.Repeat([]byte{0xff}, documentChunkSize*3)

		for _, chunk := range chunkDocumentData(data) {
			encoded, err := chunk.MarshalText()
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			if len(encoded) > qrMaxPayloadBytes {
				t.Fatalf("payload is too large: %d bytes", len(encoded))
			}

			var decoded payload.Payload
			if err := decoded.UnmarshalText(encoded); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if decoded.Index != chunk.Index || decoded.Total != chunk.Total {
				t.Fatalf("metadata mismatch: %d/%d != %d/%d", decoded.Index, decoded.Total, chunk.Index, chunk.Total)
			}
			if !bytes.Equal(decoded.Data, chunk.Data) {
				t.Fatal("decoded data mismatch")
			}
		}
	})
}
