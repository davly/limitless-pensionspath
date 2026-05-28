# CONTEXT — limitless-pensionspath

**Repo:** `C:\limitless\flagships\limitless-pensionspath`
**GitHub:** https://github.com/davly/limitless-pensionspath (public, Apache-2.0)
**Substrate:** Go 1.22
**Status:** Phase 1 scaffold (I51 marathon 2026-05-28 batch of 6 deferred NEW flagships).
**Cohort posture:** R174 5-of-5 cohort maturity FROM INCEPTION.

---

## What this flagship is

UK pension compliance scaffold pinning three regulator corpora:

1. **FCA Handbook COBS Chapter 19** (Pensions and Retirement) — amended ~monthly.
2. **TPR Defined-Benefit Funding Code of Practice** (effective September 2024).
3. **HMRC Pensions Tax Manual** — LTA abolished 2024-04-06; LSA £268,275; LSDBA £1,073,100; AA £60,000.

PLUS the cohort 5-of-5 invariants — Mirror-Mark v1 / R143 LoudOnce / R145.C firewall / R150 manifest / R151 KAT-1 anchor.

---

## R174 5-of-5 cohort maturity (strict from inception)

| Package | Discipline |
|---|---|
| `internal/firewall/` | R145.C FIREWALL-TEST-DISCIPLINE |
| `internal/honest/` | R143 LOUD-ONCE-WARNING-FLAG (5 advisories) |
| `internal/mirrormark/` | L43 Mirror-Mark v1 + R151 KAT-1 anchor |
| `internal/manifest/` | R150 PARALLEL-MAP review-metadata envelope |
| `internal/legal/` | R166 LIABILITY-FOOTER-CONST + REVIEWED-BY-COUNSEL-FALSE |

Plus 1 domain-gate package:

| Package | Domain |
|---|---|
| `internal/pension-rules/` | FCA COBS + TPR Funding Code + HMRC PTM corpus SHA pins (THE MOAT) |

---

## R143 Advisories surfaced

5 LIMITLESS_PENSIONSPATH_* advisories:

1. `LIMITLESS_PENSIONSPATH_FCA_COBS_VERSION_PIN_REQUIRED` (Error)
2. `LIMITLESS_PENSIONSPATH_TPR_FUNDING_CODE_VERSION_PIN_REQUIRED` (Error)
3. `LIMITLESS_PENSIONSPATH_HMRC_ALLOWANCES_FROZEN` (Error)
4. `LIMITLESS_PENSIONSPATH_NOT_FINANCIAL_ADVICE` (Error)
5. `LIMITLESS_PENSIONSPATH_REVIEWED_BY_COUNSEL_FALSE` (Warn)

R143.A SEVERITY-LADDER: 4 Error + 1 Warn (pension Phase-1 ladder — all three regulator pins + FCA escape are liability-bearing).

---

## Phase-2 deferred backlog (honest disclosure)

The following Phase-2+ surfaces are NOT shipped:

1. **FCA COBS 19.1 DB-to-DC transfer suitability scorer** — full pension-transfer suitability assessment (TVAS, critical yield, advice-by-PTS gating). Phase 2.
2. **TPR DB Funding Code 2024 — fast-track/bespoke valuation classifier** — full scheme-funding valuation analysis. Phase 2.
3. **HMRC Annual Allowance + tapered AA + carry-forward calculator** — Phase 2.
4. **HMRC Lump Sum Allowance / LSDBA + transitional protection engine** — Phase 2.
5. **Auto-enrolment compliance engine** (Pensions Act 2008 Chapters 1-3 of Part 1). Phase 3.
6. **Counsel review** — Phase 1 ships R166 LIABILITY-FOOTER-CONST with `ReviewedByCounsel = false`.

The Phase-1 scaffold ships ONLY: corpus-SHA pinning + R143 advisories + Mirror-Mark + R150 manifest + R145.C firewall + R166 footer + KAT-1 R151 anchor.

---

## R85 CLEAN-PARITY anchor

This CONTEXT.md is the canonical doc-comment anchor. Divergence between this status row and runtime ground truth = R85 violation.
