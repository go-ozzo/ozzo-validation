# Roadmap

## Current: v4.5.0

Restores the v4.3.0 empty-value semantics with explicit opt-ins, adds error introspection (`errors.Is`, `Flatten`)
and a build tag without `text/template`, and requires Go 1.21. See [CHANGELOG.md](CHANGELOG.md).

## Done

- [x] Performance optimization: struct field offset cache (v4.4.1)
- [x] Numeric `is.Latitude` / `is.Longitude` ([#185](https://github.com/go-ozzo/ozzo-validation/issues/185), v4.5.0)
- [x] Missing ISO 4217 currency codes VES/VED ([#206](https://github.com/go-ozzo/ozzo-validation/issues/206), v4.5.0)
- [x] `DateRule.Location()` instead of assuming UTC ([#166](https://github.com/go-ozzo/ozzo-validation/issues/166), v4.5.0)
- [x] `errors.Is` support on `Errors` ([#116](https://github.com/go-ozzo/ozzo-validation/issues/116), v4.5.0)
- [x] `Flatten` for a flat list of field errors (v4.5.0)
- [x] go.mod on Go 1.21 (v4.5.0)

Not planned for v4: `errors.As` (or `Unwrap() []error`) that looks inside `Errors`. It would make
`errors.As(err, &validation.Error)` match a map of field errors, which existing code uses to tell a single error from
a field map. Multi-error semantics are a v5 question.

## Medium Term

- [ ] `AsRule` — reuse struct validations as rules ([#167](https://github.com/go-ozzo/ozzo-validation/issues/167))
- [ ] Performance benchmarks in README

## Long Term

- [ ] Generics-based typed validation API (alongside existing API)
- [ ] TinyGo compatibility ([#163](https://github.com/go-ozzo/ozzo-validation/issues/163))
- [ ] Gradual govalidator replacement with own implementations
- [ ] v5: Go multi-error semantics for `Errors`

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Issues labeled [`good first issue`](https://github.com/go-ozzo/ozzo-validation/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) are a great starting point.
