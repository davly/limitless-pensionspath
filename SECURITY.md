# SECURITY — limitless-pensionspath

## Phase-1 scope (load-bearing)

This is a Phase-1 scaffold. The deployment MUST NOT make automated pension recommendations or scheme-funding determinations without:

1. **FCA authorisation** — FSMA 2000 s.19 General Prohibition. Carrying out a regulated activity without authorisation is a criminal offence under FSMA 2000 s.23. Hosts MUST be FCA-authorised firms.
2. **Counsel review** — `internal/legal/ReviewedByCounsel = false` is the R166 honest-default.
3. **R143 advisory acknowledgement** — the 5 LIMITLESS_PENSIONSPATH_* advisories MUST be visible to every operator.
4. **Corpus-SHA cold-verification** — before any live recommendation, the corpus pins MUST be cold-verified against regulator-published canonical (FCA COBS Handbook Notice + TPR DB Funding Code + HMRC PTM).

## R166 LIABILITY-FOOTER-CONST

The constant `internal/legal/LegalLiabilityFooter` is the FSMA escape. Every recommendation payload MUST embed it verbatim until counsel review flips `ReviewedByCounsel` to true.

## Threat model — what Phase-1 DOES NOT defend against

- Compromised signing key (no KMS integration in Phase-1).
- Compromised corpus SHA (Phase-1 placeholders sha256(CorpusID-string); Phase-2 binds to regulator-published artefacts).
- Side-channel timing attacks on base64 decode or prefix compare.
- Adversarial COBS / Funding Code / PTM parsing (no formal parser in Phase-1).
