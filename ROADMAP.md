# Roadmap

## Current: v4.4.0

Revival release after 6 years. Bug fixes, community PRs, modernized CI.

## Near Term (v4.4.x patches)

- [ ] Performance optimization (reduce reflect allocations)
- [ ] `is.Latitude` / `is.Longitude` validators ([#185](https://github.com/go-ozzo/ozzo-validation/issues/185))
- [ ] Missing ISO 4217 currency codes VES/VED ([#206](https://github.com/go-ozzo/ozzo-validation/issues/206))
- [ ] Fix DateRule UTC assumption ([#166](https://github.com/go-ozzo/ozzo-validation/issues/166))
- [ ] Fix type alias + driver.Valuer interaction ([#174](https://github.com/go-ozzo/ozzo-validation/issues/174))

## Medium Term (v4.5.0+)

- [ ] `errors.Is` / `errors.As` support on `Errors` type ([#116](https://github.com/go-ozzo/ozzo-validation/issues/116))
- [ ] `AsRule` — reuse struct validations as rules ([#167](https://github.com/go-ozzo/ozzo-validation/issues/167))
- [ ] Update go.mod to Go 1.21+
- [ ] Performance benchmarks in README

## Long Term

- [ ] Generics-based typed validation API (alongside existing API)
- [ ] TinyGo compatibility ([#163](https://github.com/go-ozzo/ozzo-validation/issues/163))
- [ ] Gradual govalidator replacement with own implementations

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Issues labeled [`good first issue`](https://github.com/go-ozzo/ozzo-validation/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) are a great starting point.
