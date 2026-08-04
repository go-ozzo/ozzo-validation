# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[4.4.0]: https://github.com/go-ozzo/ozzo-validation/compare/v4.3.0...v4.4.0
[4.3.0]: https://github.com/go-ozzo/ozzo-validation/releases/tag/v4.3.0
