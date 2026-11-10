package payload

import (
	"bytes"
	"testing"
)

func TestParseKeyshard(t *testing.T) {
	p, err := Parse([]byte("PAPERSAFE:KEYSHARD:1:10:3:3q2-7w"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	k, ok := p.(*Keyshard)
	if !ok {
		t.Fatalf("expected *Keyshard, got %T", p)
	}
	if k.Domain() != DomainKeyshard {
		t.Fatalf("unexpected domain: %d", k.Domain())
	}
	if k.ShardID != 10 {
		t.Fatalf("unexpected shard ID: %d", k.ShardID)
	}
	if k.Threshold != 3 {
		t.Fatalf("unexpected threshold: %d", k.Threshold)
	}
	if !bytes.Equal(k.Data, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatalf("unexpected data: %v", k.Data)
	}
}

func TestParseDocument(t *testing.T) {
	p, err := Parse([]byte("PAPERSAFE:DOCUMENT:1:2:5:3q2-7w"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	d, ok := p.(*Document)
	if !ok {
		t.Fatalf("expected *Document, got %T", p)
	}
	if d.Domain() != DomainDocument {
		t.Fatalf("unexpected domain: %d", d.Domain())
	}
	if d.Index != 2 || d.Total != 5 {
		t.Fatalf("unexpected chunk metadata: %d/%d", d.Index, d.Total)
	}
	if !bytes.Equal(d.Data, []byte{0xde, 0xad, 0xbe, 0xef}) {
		t.Fatalf("unexpected data: %v", d.Data)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "invalid prefix",
			input: "NOTPAPERSAFE:KEYSHARD:1:10:3:3q2-7w",
			want:  ErrInvalidPrefix,
		},
		{
			name:  "missing separator after prefix",
			input: "PAPERSAFE",
			want:  ErrInvalidPrefix,
		},
		{
			name:  "unknown domain",
			input: "PAPERSAFE:UNKNOWN:1:2:5:3q2-7w",
			want:  ErrInvalidDomain,
		},
		{
			name:  "missing domain body",
			input: "PAPERSAFE:KEYSHARD",
			want:  ErrInvalidDomain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse([]byte(tt.input)); err != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

func TestParseInvalidBase64(t *testing.T) {
	if _, err := Parse([]byte("PAPERSAFE:KEYSHARD:1:10:3:***")); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConcreteUnmarshalRejectsOtherDomain(t *testing.T) {
	var k Keyshard
	if err := k.UnmarshalText([]byte("PAPERSAFE:DOCUMENT:1:2:5:3q2-7w")); err != ErrInvalidDomain {
		t.Fatalf("expected ErrInvalidDomain, got %v", err)
	}

	var d Document
	if err := d.UnmarshalText([]byte("PAPERSAFE:KEYSHARD:1:10:3:3q2-7w")); err != ErrInvalidDomain {
		t.Fatalf("expected ErrInvalidDomain, got %v", err)
	}
}
