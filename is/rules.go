// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package is provides a list of commonly used string validation rules.
package is

import (
	"regexp"
	"unicode"

	"github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/asaskevich/govalidator.v9"
)

var (
	// Email validates if a string is an email or not.
	Email = validation.NewStringRule(govalidator.IsEmail, "email")
	// URL validates if a string is a valid URL
	URL = validation.NewStringRule(govalidator.IsURL, "url")
	// RequestURL validates if a string is a valid request URL
	RequestURL = validation.NewStringRule(govalidator.IsRequestURL, "request_url")
	// RequestURI validates if a string is a valid request URI
	RequestURI = validation.NewStringRule(govalidator.IsRequestURI, "request_uri")
	// Alpha validates if a string contains English letters only (a-zA-Z)
	Alpha = validation.NewStringRule(govalidator.IsAlpha, "alpha")
	// Digit validates if a string contains digits only (0-9)
	Digit = validation.NewStringRule(isDigit, "digit")
	// Alphanumeric validates if a string contains English letters and digits only (a-zA-Z0-9)
	Alphanumeric = validation.NewStringRule(govalidator.IsAlphanumeric, "alphanumeric")
	// UTFLetter validates if a string contains unicode letters only
	UTFLetter = validation.NewStringRule(govalidator.IsUTFLetter, "utf_letter")
	// UTFDigit validates if a string contains unicode decimal digits only
	UTFDigit = validation.NewStringRule(govalidator.IsUTFDigit, "utf_digit")
	// UTFLetterNumeric validates if a string contains unicode letters and numbers only
	UTFLetterNumeric = validation.NewStringRule(govalidator.IsUTFLetterNumeric, "utf_letter_numeric")
	// UTFNumeric validates if a string contains unicode number characters (category N) only
	UTFNumeric = validation.NewStringRule(isUTFNumeric, "utf_numeric")
	// LowerCase validates if a string contains lower case unicode letters only
	LowerCase = validation.NewStringRule(govalidator.IsLowerCase, "lower_case")
	// UpperCase validates if a string contains upper case unicode letters only
	UpperCase = validation.NewStringRule(govalidator.IsUpperCase, "upper_case")
	// Hexadecimal validates if a string is a valid hexadecimal number
	Hexadecimal = validation.NewStringRule(govalidator.IsHexadecimal, "hexadecimal")
	// HexColor validates if a string is a valid hexadecimal color code
	HexColor = validation.NewStringRule(govalidator.IsHexcolor, "hex_color")
	// RGBColor validates if a string is a valid RGB color in the form of rgb(R, G, B)
	RGBColor = validation.NewStringRule(govalidator.IsRGBcolor, "rgb_color")
	// Int validates if a string is a valid integer number
	Int = validation.NewStringRule(govalidator.IsInt, "int")
	// Float validates if a string is a floating point number
	Float = validation.NewStringRule(govalidator.IsFloat, "float")
	// UUIDv3 validates if a string is a valid version 3 UUID
	UUIDv3 = validation.NewStringRule(govalidator.IsUUIDv3, "uuid_v3")
	// UUIDv4 validates if a string is a valid version 4 UUID
	UUIDv4 = validation.NewStringRule(govalidator.IsUUIDv4, "uuid_v4")
	// UUIDv5 validates if a string is a valid version 5 UUID
	UUIDv5 = validation.NewStringRule(govalidator.IsUUIDv5, "uuid_v5")
	// UUID validates if a string is a valid UUID
	UUID = validation.NewStringRule(govalidator.IsUUID, "uuid")
	// CreditCard validates if a string is a valid credit card number
	CreditCard = validation.NewStringRule(govalidator.IsCreditCard, "credit_card")
	// ISBN10 validates if a string is an ISBN version 10
	ISBN10 = validation.NewStringRule(govalidator.IsISBN10, "isbn_10")
	// ISBN13 validates if a string is an ISBN version 13
	ISBN13 = validation.NewStringRule(govalidator.IsISBN13, "isbn_13")
	// ISBN validates if a string is an ISBN (either version 10 or 13)
	ISBN = validation.NewStringRule(isISBN, "isbn")
	// JSON validates if a string is in valid JSON format
	JSON = validation.NewStringRule(govalidator.IsJSON, "json")
	// ASCII validates if a string contains ASCII characters only
	ASCII = validation.NewStringRule(govalidator.IsASCII, "ascii")
	// PrintableASCII validates if a string contains printable ASCII characters only
	PrintableASCII = validation.NewStringRule(govalidator.IsPrintableASCII, "printable_ascii")
	// Multibyte validates if a string contains multibyte characters
	Multibyte = validation.NewStringRule(govalidator.IsMultibyte, "multibyte")
	// FullWidth validates if a string contains full-width characters
	FullWidth = validation.NewStringRule(govalidator.IsFullWidth, "full_width")
	// HalfWidth validates if a string contains half-width characters
	HalfWidth = validation.NewStringRule(govalidator.IsHalfWidth, "half_width")
	// VariableWidth validates if a string contains both full-width and half-width characters
	VariableWidth = validation.NewStringRule(govalidator.IsVariableWidth, "variable_width")
	// Base64 validates if a string is encoded in Base64
	Base64 = validation.NewStringRule(govalidator.IsBase64, "base64")
	// DataURI validates if a string is a valid base64-encoded data URI
	DataURI = validation.NewStringRule(govalidator.IsDataURI, "data_uri")
	// E164 validates if a string is a valid ISO3166 Alpha 2 country code
	E164 = validation.NewStringRule(isE164Number, "e164")
	// CountryCode2 validates if a string is a valid ISO3166 Alpha 2 country code
	CountryCode2 = validation.NewStringRule(govalidator.IsISO3166Alpha2, "country_code_2")
	// CountryCode3 validates if a string is a valid ISO3166 Alpha 3 country code
	CountryCode3 = validation.NewStringRule(govalidator.IsISO3166Alpha3, "country_code_3")
	// CurrencyCode validates if a string is a valid IsISO4217 currency code.
	CurrencyCode = validation.NewStringRule(govalidator.IsISO4217, "currency_code")
	// DialString validates if a string is a valid dial string that can be passed to Dial()
	DialString = validation.NewStringRule(govalidator.IsDialString, "dial_string")
	// MAC validates if a string is a MAC address
	MAC = validation.NewStringRule(govalidator.IsMAC, "mac")
	// IP validates if a string is a valid IP address (either version 4 or 6)
	IP = validation.NewStringRule(govalidator.IsIP, "ip")
	// IPv4 validates if a string is a valid version 4 IP address
	IPv4 = validation.NewStringRule(govalidator.IsIPv4, "ipv4")
	// IPv6 validates if a string is a valid version 6 IP address
	IPv6 = validation.NewStringRule(govalidator.IsIPv6, "ipv6")
	// Subdomain validates if a string is valid subdomain
	Subdomain = validation.NewStringRule(isSubdomain, "sub_domain")
	// Domain validates if a string is valid domain
	Domain = validation.NewStringRule(isDomain, "domain")
	// DNSName validates if a string is valid DNS name
	DNSName = validation.NewStringRule(govalidator.IsDNSName, "dns_name")
	// Host validates if a string is a valid IP (both v4 and v6) or a valid DNS name
	Host = validation.NewStringRule(govalidator.IsHost, "host")
	// Port validates if a string is a valid port number
	Port = validation.NewStringRule(govalidator.IsPort, "port")
	// MongoID validates if a string is a valid Mongo ID
	MongoID = validation.NewStringRule(govalidator.IsMongoID, "mongo_id")
	// Latitude validates if a string is a valid latitude
	Latitude = validation.NewStringRule(govalidator.IsLatitude, "latitude")
	// Longitude validates if a string is a valid longitude
	Longitude = validation.NewStringRule(govalidator.IsLongitude, "longitude")
	// SSN validates if a string is a social security number (SSN)
	SSN = validation.NewStringRule(govalidator.IsSSN, "ssn")
	// Semver validates if a string is a valid semantic version
	Semver = validation.NewStringRule(govalidator.IsSemver, "semver")
)

var (
	reDigit = regexp.MustCompile("^[0-9]+$")
	// Subdomain regex source: https://stackoverflow.com/a/7933253
	reSubdomain = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?$`)
	// E164 regex source: https://stackoverflow.com/a/23299989
	reE164 = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	// Domain regex source: https://stackoverflow.com/a/7933253
	// Slightly modified: Removed 255 max length validation since Go regex does not
	// support lookarounds. More info: https://stackoverflow.com/a/38935027
	reDomain = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z]{1,63}| xn--[a-z0-9]{1,59})$`)
)

func isISBN(value string) bool {
	return govalidator.IsISBN(value, 10) || govalidator.IsISBN(value, 13)
}

func isDigit(value string) bool {
	return reDigit.MatchString(value)
}

func isE164Number(value string) bool {
	return reE164.MatchString(value)
}

func isSubdomain(value string) bool {
	return reSubdomain.MatchString(value)
}

func isDomain(value string) bool {
	if len(value) > 255 {
		return false
	}

	return reDomain.MatchString(value)
}

func isUTFNumeric(value string) bool {
	for _, c := range value {
		if unicode.IsNumber(c) == false {
			return false
		}
	}
	return true
}
