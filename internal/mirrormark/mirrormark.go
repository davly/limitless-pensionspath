// Package mirrormark implements the cohort L43 Mirror-Mark v1 receipt
// algorithm — byte-identical to foundation/pkg/mirrormark and every
// cohort Go port.
//
// Why limitless-pensionspath consumes this today:
//
//   - Pension recommendation payloads (FCA suitability assertions, TPR
//     scheme-funding statements, HMRC annual/lifetime allowance
//     calculations) cross a regulator-grade trust boundary. A
//     recommendation payload that arrives with a verifiable Mirror-Mark
//     is provenance-anchored — an FCA / TPR / HMRC audit can cold-verify.
//   - The corpus prefix carries the pension-rules-version SHA — IS the
//     moat (FCA COBS, TPR Funding Code, HMRC Pensions Tax Manual).
package mirrormark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

const (
	MarkVersion         byte = 0x01
	MarkPrefix               = "lore@v1:"
	MarkCorpusPrefixLen      = 8
	MarkBodyLen              = MarkCorpusPrefixLen + sha256.Size
)

var (
	ErrUnknownMarkVersion = errors.New("mirrormark: unknown mark version (missing 'lore@v1:' prefix)")
	ErrMalformedMark      = errors.New("mirrormark: malformed mark (base64url decode failed or wrong body length)")
	ErrCorpusMismatch     = errors.New("mirrormark: corpus prefix mismatch (mark signed by different corpus)")
	ErrSignatureMismatch  = errors.New("mirrormark: HMAC signature mismatch (payload tampered or wrong key)")
)

func Sign(corpusSHA [sha256.Size]byte, payload []byte, key []byte) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte{MarkVersion})
	_, _ = mac.Write(corpusSHA[:])
	_, _ = mac.Write(payload)
	digest := mac.Sum(nil)
	body := make([]byte, 0, MarkBodyLen)
	body = append(body, corpusSHA[:MarkCorpusPrefixLen]...)
	body = append(body, digest...)
	return MarkPrefix + base64.RawURLEncoding.EncodeToString(body)
}

func Verify(mark string, corpusSHA [sha256.Size]byte, payload []byte, key []byte) error {
	if len(mark) < len(MarkPrefix) || mark[:len(MarkPrefix)] != MarkPrefix {
		return ErrUnknownMarkVersion
	}
	body, err := base64.RawURLEncoding.DecodeString(mark[len(MarkPrefix):])
	if err != nil {
		return ErrMalformedMark
	}
	if len(body) != MarkBodyLen {
		return ErrMalformedMark
	}
	corpusPrefix := body[:MarkCorpusPrefixLen]
	digest := body[MarkCorpusPrefixLen:]
	if !hmac.Equal(corpusPrefix, corpusSHA[:MarkCorpusPrefixLen]) {
		return ErrCorpusMismatch
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte{MarkVersion})
	_, _ = mac.Write(corpusSHA[:])
	_, _ = mac.Write(payload)
	want := mac.Sum(nil)
	if !hmac.Equal(digest, want) {
		return ErrSignatureMismatch
	}
	return nil
}
