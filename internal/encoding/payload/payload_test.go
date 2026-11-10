package payload

import (
	"bytes"
	"testing"
)

func TestMarshalTextChunkFormat(t *testing.T) {
	p := NewChunk(DomainDocument, 2, 5, []byte{0xde, 0xad, 0xbe, 0xef})

	text, err := p.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := "PAPERSAFE:DOCUMENT:2:5:3q2-7w"
	if string(text) != expected {
		t.Fatalf("text mismatch: %s != %s", text, expected)
	}
}

func TestMarshalUnmarshalChunk(t *testing.T) {
	original := NewChunk(DomainDocument, 3, 7, []byte{1, 2, 3, 4})

	text, err := original.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Payload
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if original.Domain != decoded.Domain {
		t.Fatalf("domain mismatch: %d != %d", original.Domain, decoded.Domain)
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

func TestUnmarshalTextErrors(t *testing.T) {
	var p Payload

	err := p.UnmarshalText([]byte("NOTPAPERSAFE:KEYSHARD:0a:3q2-7w"))
	if err != ErrInvalidPrefix {
		t.Fatalf("expected ErrInvalidPrefix, got %v", err)
	}

	err = p.UnmarshalText([]byte("PAPERSAFE"))
	if err != ErrInvalidPrefix {
		t.Fatalf("expected ErrInvalidPrefix, got %v", err)
	}

	err = p.UnmarshalText([]byte("PAPERSAFE:UNKNOWN:3q2-7w"))
	if err != ErrInvalidDomain {
		t.Fatalf("expected ErrInvalidDomain, got %v", err)
	}

	err = p.UnmarshalText([]byte("PAPERSAFE:KEYSHARD"))
	if err != ErrInvalidDomain {
		t.Fatalf("expected ErrInvalidDomain, got %v", err)
	}

	err = p.UnmarshalText([]byte("PAPERSAFE:KEYSHARD:0a:***"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMarshalTextInvalidChunk(t *testing.T) {
	tests := []Payload{
		{Domain: DomainDocument, Index: 1},
		{Domain: DomainDocument, Total: 1},
		{Domain: DomainDocument, Index: 2, Total: 1},
		{Domain: DomainDocument, Index: -1, Total: 1},
		{Domain: DomainDocument, Index: 1, Total: -1},
	}

	for _, p := range tests {
		if _, err := p.MarshalText(); err != ErrInvalidChunk {
			t.Fatalf("expected ErrInvalidChunk for %+v, got %v", p, err)
		}
	}
}

func TestUnmarshalTextInvalidChunk(t *testing.T) {
	inputs := []string{
		"PAPERSAFE:DOCUMENT:3q2-7w",
		"PAPERSAFE:DOCUMENT:1:3q2-7w",
		"PAPERSAFE:DOCUMENT:x:1:3q2-7w",
		"PAPERSAFE:DOCUMENT:1:x:3q2-7w",
		"PAPERSAFE:DOCUMENT:0:1:3q2-7w",
		"PAPERSAFE:DOCUMENT:1:0:3q2-7w",
		"PAPERSAFE:DOCUMENT:2:1:3q2-7w",
	}

	for _, input := range inputs {
		var p Payload
		if err := p.UnmarshalText([]byte(input)); err != ErrInvalidChunk {
			t.Fatalf("expected ErrInvalidChunk for %q, got %v", input, err)
		}
	}
}

func TestMarshalTextKeyshardFormat(t *testing.T) {
	p := NewKeyshard(0x0a, []byte{0xde, 0xad, 0xbe, 0xef})

	text, err := p.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	expected := "PAPERSAFE:KEYSHARD:0a:3q2-7w"
	if string(text) != expected {
		t.Fatalf("text mismatch: %s != %s", text, expected)
	}
}

func TestMarshalUnmarshalKeyshard(t *testing.T) {
	original := NewKeyshard(0xff, []byte{1, 2, 3, 4})

	text, err := original.MarshalText()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Payload
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if original.Domain != decoded.Domain {
		t.Fatalf("domain mismatch: %d != %d", original.Domain, decoded.Domain)
	}
	if original.ShardID != decoded.ShardID {
		t.Fatalf("shard ID mismatch: %d != %d", original.ShardID, decoded.ShardID)
	}
	if original.Index != decoded.Index || original.Total != decoded.Total {
		t.Fatalf("chunk metadata mismatch: %d/%d != %d/%d", original.Index, original.Total, decoded.Index, decoded.Total)
	}
	if !bytes.Equal(original.Data, decoded.Data) {
		t.Fatalf("data mismatch: %v != %v", original.Data, decoded.Data)
	}
}

func TestUnmarshalTextInvalidShardID(t *testing.T) {
	inputs := []string{
		"PAPERSAFE:KEYSHARD:3q2-7w",
		"PAPERSAFE:KEYSHARD:00:3q2-7w",
		"PAPERSAFE:KEYSHARD:100:3q2-7w",
		"PAPERSAFE:KEYSHARD:zz:3q2-7w",
		"PAPERSAFE:KEYSHARD::3q2-7w",
	}

	for _, input := range inputs {
		var p Payload
		if err := p.UnmarshalText([]byte(input)); err != ErrInvalidShardID {
			t.Fatalf("expected ErrInvalidShardID for %q, got %v", input, err)
		}
	}
}

func TestMarshalTextDomainMetadataValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload Payload
		want    error
	}{
		{
			name:    "keyshard with chunk metadata",
			payload: Payload{Domain: DomainKeyshard, ShardID: 1, Index: 1, Total: 1},
			want:    ErrInvalidChunk,
		},
		{
			name:    "keyshard without shard ID",
			payload: Payload{Domain: DomainKeyshard, Data: []byte{1}},
			want:    ErrInvalidShardID,
		},
		{
			name:    "document with shard ID",
			payload: Payload{Domain: DomainDocument, ShardID: 1, Index: 1, Total: 1},
			want:    ErrInvalidShardID,
		},
		{
			name:    "document without chunk metadata",
			payload: Payload{Domain: DomainDocument, Data: []byte{1}},
			want:    ErrInvalidChunk,
		},
		{
			name:    "unknown domain",
			payload: Payload{Domain: Domain(99)},
			want:    ErrInvalidDomain,
		},
	}

	for _, tt := range tests {
		if _, err := tt.payload.MarshalText(); err != tt.want {
			t.Fatalf("%s: expected %v, got %v", tt.name, tt.want, err)
		}
	}
}
