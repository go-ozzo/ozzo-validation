// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package is

// iso4217CurrencyCodes is the set of active ISO 4217 alphabetic currency codes.
//
// Source: ISO 4217, whose Registration Authority is SIX Group
// (https://www.six-group.com/en/products-services/financial-information/data-standards.html).
// Snapshot taken 2026-08-14, cross-checked against the "datasets/currency-codes" project
// (https://github.com/datasets/currency-codes), itself derived from the SIX Group lists.
//
// govalidator's bundled ISO4217List (the previous source for this rule) is missing several
// codes introduced after its last release in 2021 (e.g. MRU, SLE, VED, XAD, XCG, ZWG) and
// govalidator has had no tagged release since, so the list is maintained here instead.
//
// This list should be refreshed periodically against the ISO 4217 Maintenance Agency
// publication, since currencies are occasionally added, renamed, or withdrawn.
var iso4217CurrencyCodes = map[string]struct{}{
	"AED": {}, "AFN": {}, "ALL": {}, "AMD": {}, "AOA": {}, "ARS": {}, "AUD": {}, "AWG": {}, "AZN": {}, "BAM": {},
	"BBD": {}, "BDT": {}, "BHD": {}, "BIF": {}, "BMD": {}, "BND": {}, "BOB": {}, "BOV": {}, "BRL": {}, "BSD": {},
	"BTN": {}, "BWP": {}, "BYN": {}, "BZD": {}, "CAD": {}, "CDF": {}, "CHE": {}, "CHF": {}, "CHW": {}, "CLF": {},
	"CLP": {}, "CNY": {}, "COP": {}, "COU": {}, "CRC": {}, "CUP": {}, "CVE": {}, "CZK": {}, "DJF": {}, "DKK": {},
	"DOP": {}, "DZD": {}, "EGP": {}, "ERN": {}, "ETB": {}, "EUR": {}, "FJD": {}, "FKP": {}, "GBP": {}, "GEL": {},
	"GHS": {}, "GIP": {}, "GMD": {}, "GNF": {}, "GTQ": {}, "GYD": {}, "HKD": {}, "HNL": {}, "HTG": {}, "HUF": {},
	"IDR": {}, "ILS": {}, "INR": {}, "IQD": {}, "IRR": {}, "ISK": {}, "JMD": {}, "JOD": {}, "JPY": {}, "KES": {},
	"KGS": {}, "KHR": {}, "KMF": {}, "KPW": {}, "KRW": {}, "KWD": {}, "KYD": {}, "KZT": {}, "LAK": {}, "LBP": {},
	"LKR": {}, "LRD": {}, "LSL": {}, "LYD": {}, "MAD": {}, "MDL": {}, "MGA": {}, "MKD": {}, "MMK": {}, "MNT": {},
	"MOP": {}, "MRU": {}, "MUR": {}, "MVR": {}, "MWK": {}, "MXN": {}, "MXV": {}, "MYR": {}, "MZN": {}, "NAD": {},
	"NGN": {}, "NIO": {}, "NOK": {}, "NPR": {}, "NZD": {}, "OMR": {}, "PAB": {}, "PEN": {}, "PGK": {}, "PHP": {},
	"PKR": {}, "PLN": {}, "PYG": {}, "QAR": {}, "RON": {}, "RSD": {}, "RUB": {}, "RWF": {}, "SAR": {}, "SBD": {},
	"SCR": {}, "SDG": {}, "SEK": {}, "SGD": {}, "SHP": {}, "SLE": {}, "SOS": {}, "SRD": {}, "SSP": {}, "STN": {},
	"SVC": {}, "SYP": {}, "SZL": {}, "THB": {}, "TJS": {}, "TMT": {}, "TND": {}, "TOP": {}, "TRY": {}, "TTD": {},
	"TWD": {}, "TZS": {}, "UAH": {}, "UGX": {}, "USD": {}, "USN": {}, "UYI": {}, "UYU": {}, "UYW": {}, "UZS": {},
	"VED": {}, "VES": {}, "VND": {}, "VUV": {}, "WST": {}, "XAD": {}, "XAF": {}, "XAG": {}, "XAU": {}, "XBA": {},
	"XBB": {}, "XBC": {}, "XBD": {}, "XCD": {}, "XCG": {}, "XDR": {}, "XOF": {}, "XPD": {}, "XPF": {}, "XPT": {},
	"XSU": {}, "XTS": {}, "XUA": {}, "XXX": {}, "YER": {}, "ZAR": {}, "ZMW": {}, "ZWG": {},
}

// isISO4217CurrencyCode checks whether value is a currently active ISO 4217 alphabetic
// currency code (e.g. "USD", "EUR"). The check is case-sensitive, matching the uppercase
// form defined by the standard.
func isISO4217CurrencyCode(value string) bool {
	_, ok := iso4217CurrencyCodes[value]
	return ok
}
