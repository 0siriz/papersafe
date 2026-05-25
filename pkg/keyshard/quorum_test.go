package keyshard

import (
	"bytes"
	"testing"
)

func TestQuorumRoundtripUnsealed(t *testing.T) {
	q, err := NewQuorum()
	if err != nil {
		t.Fatalf("NewQuorum failed: %v", err)
	}

	shardSets, err := q.MakeKeyshards(5, 3)
	if err != nil {
		t.Fatalf("MakeKeyshards failed: %v", err)
	}

	reconstructed, err := ReconstructQuorum(shardSets)
	if err != nil {
		t.Fatalf("ReconstructQuorum failed: %v", err)
	}

	if !bytes.Equal(q.Key, reconstructed.Key) {
		t.Fatal("Key mismatch")
	}

	if !bytes.Equal(q.SigningKey, reconstructed.SigningKey) {
		t.Fatal("SigningKey mismatch")
	}

	if reconstructed.IsSealed() {
		t.Fatal("reconstructed unsealed quorum should not be sealed")
	}

	_, err = reconstructed.MakeKeyshards(3, 2)
	if err != nil {
		t.Fatal("reconstructed unsealed quorum should be able to issue shards")
	}
}

func TestQuorumRoundtripSealed(t *testing.T) {
	q, err := NewQuorum()
	if err != nil {
		t.Fatalf("NewQuorum failed: %v", err)
	}

	q.SetSealed()
	if !q.IsSealed() {
		t.Fatal("quorum should be sealed after SetSealed")
	}

	shardSets, err := q.MakeKeyshards(5, 3)
	if err != nil {
		t.Fatalf("MakeKeyshards failed: %v", err)
	}

	reconstructed, err := ReconstructQuorum(shardSets)
	if err != nil {
		t.Fatalf("ReconstructQuorum failed: %v", err)
	}

	if !bytes.Equal(q.Key, reconstructed.Key) {
		t.Fatal("Key mismatch")
	}

	if len(reconstructed.SigningKey) != 0 {
		t.Fatal("sealed reconstruction should not recover signing key")
	}

	if !reconstructed.IsSealed() {
		t.Fatal("reconstructed sealed quorum should be sealed")
	}
}

func TestQuorumDestroy(t *testing.T) {
	q, err := NewQuorum()
	if err != nil {
		t.Fatalf("NewQuorum failed: %v", err)
	}

	q.Destroy()
	if !q.IsSealed() {
		t.Fatal("quorum should be sealed after Destroy")
	}

	expected := make([]byte, len(q.SigningKey))
	if !bytes.Equal(q.SigningKey, expected) {
		t.Fatal("signing key should be zeroed after Destroy")
	}
}

func TestQuorumMultipleIssuesUnsealed(t *testing.T) {
	q, err := NewQuorum()
	if err != nil {
		t.Fatalf("NewQuorum failed: %v", err)
	}

	sets1, err := q.MakeKeyshards(5, 3)
	if err != nil {
		t.Fatalf("first MakeKeyshards failed: %v", err)
	}

	sets2, err := q.MakeKeyshards(3, 2)
	if err != nil {
		t.Fatalf("second MakeKeyshards failed: %v", err)
	}

	r1, _ := ReconstructQuorum(sets1)
	r2, _ := ReconstructQuorum(sets2)

	if !bytes.Equal(r1.Key, r2.Key) {
		t.Fatal("both reconstructions should recover the same key")
	}

	if !bytes.Equal(r1.SigningKey, r2.SigningKey) {
		t.Fatal("both reconstructions should recover the same signing key")
	}
}
