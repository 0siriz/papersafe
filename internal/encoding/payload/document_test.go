package payload

import (
	"bytes"
	"testing"
)

func TestMarshalTextDocumentFormat(t *testing.T) {
	p := NewDocument(2, 5, []byte{0xde, 0xad, 0xbe, 0xef})

	text, err := p.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := "PAPERSAFE:DOCUMENT:1:2:5:3q2-7w"
	if string(text) != expected {
		t.Fatalf("text mismatch: %s != %s", text, expected)
	}
}

func TestMarshalUnmarshalDocument(t *testing.T) {
	original := NewDocument(3, 7, []byte{1, 2, 3, 4})

	text, err := original.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Document
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if original.Index != decoded.Index {
		t.Fatalf("index mismatch: %d != %d", original.Index, decoded.Index)
	}
	if original.Total != decoded.Total {
		t.Fatalf("total mismatch: %d != %d", original.Total, decoded.Total)
	}
	if !bytes.Equal(original.Data, decoded.Data) {
		t.Fatalf("data mismatch: %v != %v", original.Data, decoded.Data)
	}
}

func TestMarshalTextInvalidDocument(t *testing.T) {
	tests := []struct {
		name    string
		payload Document
		want    error
	}{
		{
			name:    "missing index",
			payload: NewDocument(0, 1, []byte{1}),
			want:    ErrInvalidChunk,
		},
		{
			name:    "missing total",
			payload: NewDocument(1, 0, []byte{1}),
			want:    ErrInvalidChunk,
		},
		{
			name:    "index exceeds total",
			payload: NewDocument(2, 1, []byte{1}),
			want:    ErrInvalidChunk,
		},
		{
			name:    "negative index",
			payload: NewDocument(-1, 1, []byte{1}),
			want:    ErrInvalidChunk,
		},
		{
			name:    "negative total",
			payload: NewDocument(1, -1, []byte{1}),
			want:    ErrInvalidChunk,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.payload.MarshalText(); err != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestUnmarshalTextInvalidDocument(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "missing version",
			input: "PAPERSAFE:DOCUMENT::2:5:3q2-7w",
			want:  ErrInvalidVersion,
		},
		{
			name:  "non-numeric version",
			input: "PAPERSAFE:DOCUMENT:x:2:5:3q2-7w",
			want:  ErrInvalidVersion,
		},
		{
			name:  "unsupported version",
			input: "PAPERSAFE:DOCUMENT:2:2:5:3q2-7w",
			want:  ErrUnsupportedVersion,
		},
		{
			name:  "missing index",
			input: "PAPERSAFE:DOCUMENT:1",
			want:  ErrInvalidVersion,
		},
		{
			name:  "missing total",
			input: "PAPERSAFE:DOCUMENT:1:2",
			want:  ErrInvalidChunk,
		},
		{
			name:  "missing data",
			input: "PAPERSAFE:DOCUMENT:1:2:5",
			want:  ErrInvalidChunk,
		},
		{
			name:  "non-numeric index",
			input: "PAPERSAFE:DOCUMENT:1:x:5:3q2-7w",
			want:  ErrInvalidChunk,
		},
		{
			name:  "non-numeric total",
			input: "PAPERSAFE:DOCUMENT:1:2:x:3q2-7w",
			want:  ErrInvalidChunk,
		},
		{
			name:  "zero index",
			input: "PAPERSAFE:DOCUMENT:1:0:5:3q2-7w",
			want:  ErrInvalidChunk,
		},
		{
			name:  "zero total",
			input: "PAPERSAFE:DOCUMENT:1:2:0:3q2-7w",
			want:  ErrInvalidChunk,
		},
		{
			name:  "index exceeds total",
			input: "PAPERSAFE:DOCUMENT:1:6:5:3q2-7w",
			want:  ErrInvalidChunk,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Document
			if err := p.UnmarshalText([]byte(tt.input)); err != tt.want {
				t.Fatalf("expected %v for %q, got %v", tt.want, tt.input, err)
			}
		})
	}
}

func TestUnmarshalTextDocumentInvalidBase64(t *testing.T) {
	var p Document
	if err := p.UnmarshalText([]byte("PAPERSAFE:DOCUMENT:1:2:5:***")); err == nil {
		t.Fatal("expected error, got nil")
	}
}
