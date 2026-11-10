package payload

import (
	"bytes"
	"testing"
)

func TestMarshalTextKeyshardFormat(t *testing.T) {
	p := NewKeyshard(10, 3, []byte{0xde, 0xad, 0xbe, 0xef})

	text, err := p.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := "PAPERSAFE:KEYSHARD:1:10:3:3q2-7w"
	if string(text) != expected {
		t.Fatalf("text mismatch: %s != %s", text, expected)
	}
}

func TestMarshalUnmarshalKeyshard(t *testing.T) {
	original := NewKeyshard(255, 3, []byte{1, 2, 3, 4})

	text, err := original.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Keyshard
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if original.ShardID != decoded.ShardID {
		t.Fatalf("shard ID mismatch: %d != %d", original.ShardID, decoded.ShardID)
	}
	if original.Threshold != decoded.Threshold {
		t.Fatalf("threshold mismatch: %d != %d", original.Threshold, decoded.Threshold)
	}
	if !bytes.Equal(original.Data, decoded.Data) {
		t.Fatalf("data mismatch: %v != %v", original.Data, decoded.Data)
	}
}

func TestMarshalTextInvalidKeyshard(t *testing.T) {
	tests := []struct {
		name    string
		payload Keyshard
		want    error
	}{
		{
			name:    "missing shard ID",
			payload: NewKeyshard(0, 3, []byte{1}),
			want:    ErrInvalidShardID,
		},
		{
			name:    "threshold too small",
			payload: NewKeyshard(1, 1, []byte{1}),
			want:    ErrInvalidShardParams,
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

func TestUnmarshalTextInvalidKeyshard(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "missing version",
			input: "PAPERSAFE:KEYSHARD::10:3:3q2-7w",
			want:  ErrInvalidVersion,
		},
		{
			name:  "non-numeric version",
			input: "PAPERSAFE:KEYSHARD:x:10:3:3q2-7w",
			want:  ErrInvalidVersion,
		},
		{
			name:  "unsupported version",
			input: "PAPERSAFE:KEYSHARD:2:10:3:3q2-7w",
			want:  ErrUnsupportedVersion,
		},
		{
			name:  "missing shard ID",
			input: "PAPERSAFE:KEYSHARD:1",
			want:  ErrInvalidVersion,
		},
		{
			name:  "zero shard ID",
			input: "PAPERSAFE:KEYSHARD:1:0:3:3q2-7w",
			want:  ErrInvalidShardID,
		},
		{
			name:  "shard ID too large",
			input: "PAPERSAFE:KEYSHARD:1:256:3:3q2-7w",
			want:  ErrInvalidShardID,
		},
		{
			name:  "non-numeric shard ID",
			input: "PAPERSAFE:KEYSHARD:1:zz:3:3q2-7w",
			want:  ErrInvalidShardID,
		},
		{
			name:  "empty shard ID",
			input: "PAPERSAFE:KEYSHARD:1::3:3q2-7w",
			want:  ErrInvalidShardID,
		},
		{
			name:  "missing threshold",
			input: "PAPERSAFE:KEYSHARD:1:10",
			want:  ErrInvalidShardID,
		},
		{
			name:  "missing data",
			input: "PAPERSAFE:KEYSHARD:1:10:3",
			want:  ErrInvalidShardParams,
		},
		{
			name:  "empty threshold",
			input: "PAPERSAFE:KEYSHARD:1:10::3q2-7w",
			want:  ErrInvalidShardParams,
		},
		{
			name:  "non-numeric threshold",
			input: "PAPERSAFE:KEYSHARD:1:10:x:3q2-7w",
			want:  ErrInvalidShardParams,
		},
		{
			name:  "threshold too small",
			input: "PAPERSAFE:KEYSHARD:1:10:1:3q2-7w",
			want:  ErrInvalidShardParams,
		},
		{
			name:  "threshold too large",
			input: "PAPERSAFE:KEYSHARD:1:10:256:3q2-7w",
			want:  ErrInvalidShardParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p Keyshard
			if err := p.UnmarshalText([]byte(tt.input)); err != tt.want {
				t.Fatalf("expected %v for %q, got %v", tt.want, tt.input, err)
			}
		})
	}
}

func TestUnmarshalTextKeyshardInvalidBase64(t *testing.T) {
	var p Keyshard
	if err := p.UnmarshalText([]byte("PAPERSAFE:KEYSHARD:1:10:3:***")); err == nil {
		t.Fatal("expected error, got nil")
	}
}
