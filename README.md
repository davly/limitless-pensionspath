# limitless-pensionspath

**One-line:** UK FCA + TPR + HMRC pension rules compliance forge — corpus-pinned methodology with R166 FSMA 2000 s.19 liability-footer escape.

**Category:** B2B Enterprise | RegTech / Pension Compliance
**Target Market:** FCA-authorised independent financial advisers (IFAs), pension transfer specialists (PTS), pension administration firms, occupational pension scheme trustees, scheme actuaries.
**Trinity Engines:** Reality (primary — corpus-pinned rules), Causal (secondary — outcome attribution), Parallax (secondary — cross-regime contradiction).

**Status:** Phase-1 scaffold (2026-05-28 marathon I51 batch). R174 5-of-5 cohort maturity FROM INCEPTION.

---

## Problem Statement

UK pension regulation spans three independently-moving regulators:

- **FCA** publishes COBS Handbook updates ~monthly via Handbook Notices. COBS 19 (Pensions + Retirement), particularly COBS 19.1 (DB-to-DC transfer advice), is the most-amended retail conduct chapter.
- **TPR** issued the new DB Funding Code of Practice effective September 2024 under Pension Schemes Act 2021 + Occupational Pension Schemes (Funding and Investment Strategy and Amendment) Regulations 2024 (SI 2024/462) — load-bearing for every DB scheme valuation.
- **HMRC** Pensions Tax Manual reflects rapidly-shifting tax thresholds: Lifetime Allowance abolished 2024-04-06 (replaced by LSA £268,275 + LSDBA £1,073,100); Annual Allowance now £60,000 (was £40,000 pre-2023-04-06); Money Purchase AA £10,000.

Every pension-rules compliance vendor either (a) implements ONE regulator's rules in isolation, OR (b) ships hard-coded thresholds that go stale within months. The fundamental problem: an FCA s.166 review, TPR scheme-funding investigation, or HMRC IHTM/PTM audit cannot answer *which rule version was used, and is it still current?*

---

## R166 FSMA 2000 Liability-Footer (load-bearing)

Every recommendation emitted by this software carries the R166 LIABILITY-FOOTER-CONST escape:

> **NOT FINANCIAL ADVICE.** Regulated financial advice in the UK is restricted by FSMA 2000 s.19 (the General Prohibition) to FCA-authorised firms. Carrying out a regulated activity without authorisation is a criminal offence under FSMA 2000 s.23. This software ships a Phase-1 scaffold and is NOT a substitute for FCA-authorised independent financial advice or actuarial review.

This is the FSMA firewall. The library `LibraryRecommendsHostActs = false` declares that downstream hosts MUST surface the footer before acting on any recommendation.

---

## R174 5-of-5 cohort maturity (strict from inception)

- **L43 Mirror-Mark v1** — `internal/mirrormark/`. KAT-1 hex `239a7d0d3f1bbe3a98aede01e2ad818c2db60b7177c02e2f015035b2b5b7dbca` pinned + OpenSSL-reproducible.
- **R143 LOUD-ONCE-WARNING-FLAG** — `internal/honest/` with 5 LIMITLESS_PENSIONSPATH_* advisories (4 Error + 1 Warn).
- **R145.C FIREWALL-TEST-DISCIPLINE** — `internal/firewall/` with on-disk drift detection.
- **R150 PARALLEL-MAP review-metadata** — `internal/manifest/` with FreshAt + Source + SchemaVersion + Confidence + Jurisdiction + Version (Class-3 anchor).
- **R151 KAT-AS-COHORT-INVARIANT-PIN** — KAT-1 byte-identical to every cohort substrate.

---

## Phase-2 deferred backlog

See CONTEXT.md "Phase-2 deferred backlog" section.

---

## License

Apache-2.0.
