# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.5.0] - 2026-10-27

This release restores the v4.3.0 behavior that v4.4.0 changed for empty values. It adds explicit opt-ins for the
cases v4.4.0 was meant to fix, and fixes the problems found in the unreleased changes before they ship.

### Fixed
- Restore the v4.3.0 empty-value semantics ([#234](https://github.com/go-ozzo/ozzo-validation/issues/234), reported
  by [@ldez](https://github.com/ldez)). Two "Breaking Changes" of v4.4.0 are reverted:
  - a zero struct is no longer empty: only the zero `time.Time` is. This applies to every rule that checks for
    empty values. `Required` and `NilOrNotEmpty` pass again for a zero struct, and `Required` on a zero
    `Validatable` struct reports the struct's own field errors instead of "cannot be blank". `Empty` reports a zero
    struct as not empty again, and the rules that skip empty values (`In`, `NotIn`, `Length`, `RuneLength`, `Date`
    and the string rules, including `is.*`) validate a zero struct again. `Each(Required)` passes again for
    zero-struct elements;
  - `Min()`/`Max()` skip empty values again (0, 0.0, the zero `time.Time`), as their documentation says.

  To keep the v4.4 checks, use the new opt-ins below.
- `Errors.Error()` no longer starts with `"; "` when the first sorted key holds a nil error. An `Errors` with only
  nil entries returns `""`, as an empty `Errors` does.
- `is.Latitude`/`is.Longitude` are `validation.StringRule` again, as in v4.4.1, so code that holds them as a
  StringRule compiles. String input is validated exactly as in v4.4.1: `"NaN"`, `"1e1"`, `"0x1p-2"`, `"45."` and
  `".5"` are rejected. Numbers are accepted as introduced by
  [#223](https://github.com/go-ozzo/ozzo-validation/pull/223); NaN and infinities are rejected.
- A missing map key never runs its rules with a nil value: rules written for the key's type, such as a `By` with a
  type assertion, no longer panic. The error stays "required key is missing" (`validation_key_missing`). Only a
  `Required` or `NotNil` rule with a custom message (`Error`/`ErrorObject`) reports its own error
  ([#222](https://github.com/go-ozzo/ozzo-validation/pull/222) by
  [@nguyenvantuan2391996](https://github.com/nguyenvantuan2391996), narrowed) — fixes
  [#134](https://github.com/go-ozzo/ozzo-validation/issues/134).
- `EnsureString(nil)` returns an error instead of panicking
  ([#232](https://github.com/go-ozzo/ozzo-validation/pull/232) by [@yuefdev](https://github.com/yuefdev)) — fixes
  [#233](https://github.com/go-ozzo/ozzo-validation/issues/233).
- `MultipleOf(0)` returns an error instead of panicking with a division by zero
  ([#231](https://github.com/go-ozzo/ozzo-validation/pull/231) by [@SashaMIT](https://github.com/SashaMIT));
  `MultipleOf(nil)` returns an error instead of panicking.
- `NotIn()` compares with `reflect.DeepEqual`, as `In()` does: it no longer panics on uncomparable values such as
  `[]byte`, and structs with pointer fields compare by the pointed-to values
  ([#229](https://github.com/go-ozzo/ozzo-validation/pull/229) by [@youdie006](https://github.com/youdie006)).

### Added
- `ozzo_validation_no_template` build tag: error messages are rendered without `text/template`, which removes it
  from the binary and makes it noticeably smaller. Built-in messages read the same; custom messages support only
  field placeholders such as `{{.min}}` and `{{ .min }}`. See the README
  ([#235](https://github.com/go-ozzo/ozzo-validation/pull/235) by [@ldez](https://github.com/ldez)) — fixes
  [#230](https://github.com/go-ozzo/ozzo-validation/issues/230).
- `ThresholdRule.ValidateEmpty()`: `Min`/`Max` also check empty values, e.g. `Min(0).Exclusive().ValidateEmpty()`
  rejects 0 ([#165](https://github.com/go-ozzo/ozzo-validation/issues/165),
  [#180](https://github.com/go-ozzo/ozzo-validation/issues/180)).
- `EmptyChecker` interface: a struct type defines when its value is empty
  ([#143](https://github.com/go-ozzo/ozzo-validation/issues/143)).
- `StringRule.AllowNumber()`: a string rule can also accept finite numbers.
- `is.Latitude`/`is.Longitude` accept numbers (int, uint and float types)
  ([#223](https://github.com/go-ozzo/ozzo-validation/pull/223) by
  [@nguyenvantuan2391996](https://github.com/nguyenvantuan2391996)) — fixes
  [#185](https://github.com/go-ozzo/ozzo-validation/issues/185).
- `DateRule.Location()`: parse dates that carry no time zone in a given location; the default stays UTC
  ([#220](https://github.com/go-ozzo/ozzo-validation/pull/220) by
  [@nguyenvantuan2391996](https://github.com/nguyenvantuan2391996)) — fixes
  [#166](https://github.com/go-ozzo/ozzo-validation/issues/166).
- `Errors.Is`: `errors.Is` finds a matching error anywhere in `validation.Errors`, including nested maps, slices and
  structs ([#225](https://github.com/go-ozzo/ozzo-validation/pull/225) by
  [@ChrisJr404](https://github.com/ChrisJr404)) — fixes [#116](https://github.com/go-ozzo/ozzo-validation/issues/116).

  #225 also proposed `Errors.As`. It is not part of this release. With it, `errors.As(err, &validation.Error)`
  returned the first field's error for a map of field errors, and code that uses that call to tell a single error from
  a field map lost the field names. `errors.As` keeps reporting `false` for `validation.Errors`; use `errors.Is` or
  `Flatten()` to reach the individual errors.
- `ErrorObject.Is`: a validation error matches a built-in error with the same code, so
  `errors.Is(err, validation.ErrRequired)` works, also for rules with a custom message such as
  `Required.Error("...")`. Give the sentinel errors of your own rules a code: an error with an empty code never
  matches.
- `Flatten(err)`, `Errors.Flatten()` and `FieldError`: a flat, sorted list of all field errors with the key path,
  code, rendered message and params, and the error itself in `Err`. `Field()` gives `items.0.name`, `Pointer()` gives
  the RFC 6901 pointer `/items/0/name`. Slice indexes sort numerically (`items.2` before `items.10`). A single error
  that is not an `Errors` is one entry with no path. Check `Err` for an `InternalError` before sending `Message` to a
  client.
- `errors.Is` and `errors.As` reach the error wrapped by `NewInternalError`.
- README: "Empty Values", "Checking for an Error", "Listing Field Errors" and "Smaller Binaries" sections.

### Changed
- Go 1.21 or later is required (`go 1.21` in go.mod), the minimum announced in
  [#207](https://github.com/go-ozzo/ozzo-validation/issues/207).
- `is.CurrencyCode` uses an ISO 4217 list maintained in this module instead of govalidator's, which adds current
  codes such as VES and VED ([#224](https://github.com/go-ozzo/ozzo-validation/pull/224) by
  [@nguyenvantuan2391996](https://github.com/nguyenvantuan2391996)) — fixes
  [#206](https://github.com/go-ozzo/ozzo-validation/issues/206).
- testify v1.4.0 → v1.12.1 (test-only).
- CI tests Go 1.21, oldstable and stable on Linux, and stable on Windows and macOS, with the race detector, with and
  without the `ozzo_validation_no_template` tag. It also builds and tests on 386, checks formatting and
  `go mod tidy`, and pins every tool version except govulncheck, which stays on the latest release by design so that
  the vulnerability data is current. One required check, "Test", gates all jobs.
- Docs: the v4.4.0 behavior changes documented
  ([#228](https://github.com/go-ozzo/ozzo-validation/pull/228) by [@dw-kihara](https://github.com/dw-kihara));
  updated in this release for the restored semantics.
- Tests: coverage gaps closed ([#218](https://github.com/go-ozzo/ozzo-validation/pull/218) by
  [@nguyenvantuan2391996](https://github.com/nguyenvantuan2391996)).
- The committed `coverage.out` is removed; test outputs are ignored.

## [4.4.1] - 2026-08-04

### Performance
- Struct field offset cache: `ValidateStruct` with 6 fields makes 36% fewer allocations, uses 50% less memory and
  runs 46% faster; `ValidateWithConditional` makes 44% fewer allocations and runs 56% faster.
- The benchmark pre-compiles its regex (it inflated the allocation count).

## [4.4.0] - 2026-08-04

Project revived after 6 years of inactivity. New maintainer: [@kolkov](https://github.com/kolkov).

### Added
- `is.UUIDv7` validation rule per RFC 9562 ([#205](https://github.com/go-ozzo/ozzo-validation/pull/205) by [@dmitryzhvinklis](https://github.com/dmitryzhvinklis))
- `is.ULID` validation rule ([#171](https://github.com/go-ozzo/ozzo-validation/pull/171) by [@upamune](https://github.com/upamune))
- `is.Origin` validation rule for CORS origins ([#198](https://github.com/go-ozzo/ozzo-validation/pull/198) by [@nguyenvantuan2391996](https://github.com/nguyenvantuan2391996))
- `OptionalKey()` and `DynamicMap()` convenience constructors ([#151](https://github.com/go-ozzo/ozzo-validation/pull/151) by [@Jessinra](https://github.com/Jessinra))
- `EachUntilFirstError()` rule for large collections ([#93](https://github.com/go-ozzo/ozzo-validation/pull/93) by [@geekflyer](https://github.com/geekflyer))
- GitHub Actions CI with Go 1.22/1.23 test matrix
- golangci-lint integration with enterprise config
- Codecov OIDC integration (99.5% coverage)
- Benchmark infrastructure (PR comparison + historical tracking)
- CONTRIBUTING.md, CODE_OF_CONDUCT.md, SECURITY.md, CODEOWNERS

### Fixed
- E.164 phone validation now requires `+` prefix per ITU-T standard ([#208](https://github.com/go-ozzo/ozzo-validation/pull/208) by [@Yanhu007](https://github.com/Yanhu007)) — fixes [#195](https://github.com/go-ozzo/ozzo-validation/issues/195)
- Domain regex allows uppercase in final character ([#197](https://github.com/go-ozzo/ozzo-validation/pull/197) by [@mhargrove](https://github.com/mhargrove))
- `Each()` now passes pointers correctly to `By()` callbacks ([#160](https://github.com/go-ozzo/ozzo-validation/pull/160) by [@dane](https://github.com/dane))
- `IsEmpty()` treats all zero-valued structs as empty, not just `time.Time{}` ([#144](https://github.com/go-ozzo/ozzo-validation/pull/144) by [@soranoba](https://github.com/soranoba)) — fixes [#143](https://github.com/go-ozzo/ozzo-validation/issues/143)
- `Errors.Error()` no longer panics on nil values ([#212](https://github.com/go-ozzo/ozzo-validation/pull/212)) — fixes [#147](https://github.com/go-ozzo/ozzo-validation/issues/147)
- `Min()`/`Max()` now validate zero values instead of silently skipping them ([#212](https://github.com/go-ozzo/ozzo-validation/pull/212)) — fixes [#165](https://github.com/go-ozzo/ozzo-validation/issues/165), [#180](https://github.com/go-ozzo/ozzo-validation/issues/180)
- README install command corrected to v4 module path ([#194](https://github.com/go-ozzo/ozzo-validation/pull/194) by [@apuatcfbd](https://github.com/apuatcfbd))
- README code examples: single quotes → double quotes ([#178](https://github.com/go-ozzo/ozzo-validation/pull/178) by [@Hannes-tallied](https://github.com/Hannes-tallied))
- README city validation example uses realistic Length(1, 50) ([#175](https://github.com/go-ozzo/ozzo-validation/pull/175) by [@tuan-nxcr](https://github.com/tuan-nxcr))
- Code formatting for Go 1.19+ doc comment style
- Various golangci-lint issues (whitespace, bool comparison, if-else chains)

### Changed
- Badges: Travis CI → GitHub Actions, Coveralls → Codecov
- govalidator dependency updated
- Removed dead Travis CI configuration

### Breaking Changes
- `is.E164` now **rejects** phone numbers without `+` prefix (spec-correct behavior)
- `Each()` with `By()` on pointer slices now passes `*T` instead of `T` (correct behavior)
- `IsEmpty()` returns `true` for zero-valued structs (was only `time.Time{}` before)
- `Min()`/`Max()` now validate numeric zero and `time.Time{}` instead of skipping them

## [4.3.0] - 2020-10-19

_Last release by original author [@qiangxue](https://github.com/qiangxue)._

[4.5.0]: https://github.com/go-ozzo/ozzo-validation/compare/v4.4.1...v4.5.0
[4.4.1]: https://github.com/go-ozzo/ozzo-validation/compare/v4.4.0...v4.4.1
[4.4.0]: https://github.com/go-ozzo/ozzo-validation/compare/v4.3.0...v4.4.0
[4.3.0]: https://github.com/go-ozzo/ozzo-validation/releases/tag/v4.3.0
