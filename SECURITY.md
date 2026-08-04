# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 4.3.x   | :white_check_mark: |

## Reporting a Vulnerability

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, please report security issues via:

1. **Private Security Advisory** (preferred):
   https://github.com/go-ozzo/ozzo-validation/security/advisories/new

2. **GitHub Issues** (for less critical issues):
   https://github.com/go-ozzo/ozzo-validation/issues

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Affected versions
- Potential impact

### Response Timeline

- **Initial Response**: Within 72 hours
- **Fix & Disclosure**: Coordinated with reporter

## Security Considerations

ozzo-validation is used at input boundaries. Users should be aware of:

1. **Regex Patterns** — Custom `Match()` rules with untrusted regex patterns can cause ReDoS
2. **Error Messages** — Validation errors may contain user input; sanitize before rendering in HTML
3. **Type Assertions** — Custom `By()` validators should handle unexpected types gracefully

## Security Contact

- **GitHub Security Advisory**: https://github.com/go-ozzo/ozzo-validation/security/advisories/new
- **Public Issues**: https://github.com/go-ozzo/ozzo-validation/issues
