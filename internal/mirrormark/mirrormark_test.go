package mirrormark

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

const kat1Mark = "lore@v1:AAAAAAAAAAAjmn0NPxu-Opiu3gHirYGMLbYLcXfALi8BUDWytbfbyg"
const kat1DigestHex = "239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca"

func TestVerify_KAT1Mark(t *testing.T) {
	var zeroCorpus [sha256.Size]byte
	if err := Verify(kat1Mark, zeroCorpus, []byte{}, []byte{}); err != nil {
		t.Fatalf("KAT-1 cohort literal failed Verify: %v", err)
	}
}

func TestSign_ProducesKAT1Mark(t *testing.T) {
	var zeroCorpus [sha256.Size]byte
	got := Sign(zeroCorpus, []byte{}, []byte{})
	if got != kat1Mark {
		t.Fatalf("KAT-1 drift: got %q want %q", got, kat1Mark)
	}
}

func TestKAT1Digest_EmbeddedInKAT1Mark(t *testing.T) {
	encoded := strings.TrimPrefix(kat1Mark, MarkPrefix)
	body, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64: %v", err)
	}
	if len(body) != MarkBodyLen {
		t.Fatalf("body len: %d", len(body))
	}
	got := hex.EncodeToString(body[MarkCorpusPrefixLen:])
	if got != kat1DigestHex {
		t.Fatalf("digest drift: got %s want %s", got, kat1DigestHex)
	}
}

func TestSign_RoundtripVerify(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = byte(i)
	}
	key := []byte("pensions_test_key")
	payload := []byte(`{"member":"NIN-AB123456A","scheme":"DC"}`)
	mark := Sign(corpus, payload, key)
	if err := Verify(mark, corpus, payload, key); err != nil {
		t.Fatalf("verify rejected: %v", err)
	}
}

func TestVerify_RejectsMissingPrefix(t *testing.T) {
	var corpus [sha256.Size]byte
	if err := Verify("not-a-mark", corpus, []byte{}, []byte("k")); err != ErrUnknownMarkVersion {
		t.Fatalf("got %v", err)
	}
}

func TestVerify_RejectsMalformedBase64(t *testing.T) {
	var corpus [sha256.Size]byte
	if err := Verify("lore@v1:!!!", corpus, []byte{}, []byte("k")); err != ErrMalformedMark {
		t.Fatalf("got %v", err)
	}
}

func TestVerify_RejectsWrongCorpus(t *testing.T) {
	var a, b [sha256.Size]byte
	for i := range a {
		a[i] = 0x11
		b[i] = 0x22
	}
	mark := Sign(a, []byte("p"), []byte("k"))
	if err := Verify(mark, b, []byte("p"), []byte("k")); err != ErrCorpusMismatch {
		t.Fatalf("got %v", err)
	}
}

func TestVerify_RejectsTamperedPayload(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = 0x44
	}
	mark := Sign(corpus, []byte("orig"), []byte("k"))
	if err := Verify(mark, corpus, []byte("tamp"), []byte("k")); err != ErrSignatureMismatch {
		t.Fatalf("got %v", err)
	}
}

func TestVerify_RejectsTamperedKey(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = 0x55
	}
	mark := Sign(corpus, []byte("p"), []byte("alice"))
	if err := Verify(mark, corpus, []byte("p"), []byte("bob")); err != ErrSignatureMismatch {
		t.Fatalf("got %v", err)
	}
}

func TestMarkLength_FixedAt62(t *testing.T) {
	var corpus [sha256.Size]byte
	for i := range corpus {
		corpus[i] = byte(i * 3)
	}
	if got := len(Sign(corpus, []byte("x"), []byte("k"))); got != 62 {
		t.Fatalf("got %d, want 62", got)
	}
}

func TestMarkPrefix_Pinned(t *testing.T) {
	if MarkPrefix != "lore@v1:" {
		t.Fatalf("drift: %q", MarkPrefix)
	}
}
