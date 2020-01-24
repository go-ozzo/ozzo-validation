// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package is provides a list of commonly used string validation rules.
package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
	"unicode"

	"gopkg.in/asaskevich/govalidator.v9"
)

var (
	// Email validates if a string is an email or not.
	Email = validation.NewStringValidator(govalidator.IsEmail, "validation_email_invalid", "must be a valid email address")
	// URL validates if a string is a valid URL
	URL = validation.NewStringValidator(govalidator.IsURL, "validation_url_invalid", "must be a valid URL")
	// RequestURL validates if a string is a valid request URL
	RequestURL = validation.NewStringValidator(govalidator.IsRequestURL, "validation_request_url_invalid", "must be a valid request URL")
	// RequestURI validates if a string is a valid request URI
	RequestURI = validation.NewStringValidator(govalidator.IsRequestURI, "validation_request_uri_invalid", "must be a valid request URI")
	// Alpha validates if a string contains English letters only (a-zA-Z)
	Alpha = validation.NewStringValidator(govalidator.IsAlpha, "validation_alpha_invalid", "must contain English letters only")
	// Digit validates if a string contains digits only (0-9)
	Digit = validation.NewStringValidator(isDigit, "validation_digit_invalid", "must contain digits only")
	// Alphanumeric validates if a string contains English letters and digits only (a-zA-Z0-9)
	Alphanumeric = validation.NewStringValidator(govalidator.IsAlphanumeric, "validation_alphanumeric_invalid", "must contain English letters and digits only")
	// UTFLetter validates if a string contains unicode letters only
	UTFLetter = validation.NewStringValidator(govalidator.IsUTFLetter, "validation_utf_letter_invalid", "must contain unicode letter characters only")
	// UTFDigit validates if a string contains unicode decimal digits only
	UTFDigit = validation.NewStringValidator(govalidator.IsUTFDigit, "validation_utf_digit_invalid", "must contain unicode decimal digits only")
	// UTFLetterNumeric validates if a string contains unicode letters and numbers only
	UTFLetterNumeric = validation.NewStringValidator(govalidator.IsUTFLetterNumeric, "validation_utf_letter_numeric_invalid", "must contain unicode letters and numbers only")
	// UTFNumeric validates if a string contains unicode number characters (category N) only
	UTFNumeric = validation.NewStringValidator(isUTFNumeric, "validation_utf_numeric_invalid", "must contain unicode number characters only")
	// LowerCase validates if a string contains lower case unicode letters only
	LowerCase = validation.NewStringValidator(govalidator.IsLowerCase, "validation_lower_case_invalid", "must be in lower case")
	// UpperCase validates if a string contains upper case unicode letters only
	UpperCase = validation.NewStringValidator(govalidator.IsUpperCase, "validation_upper_case_invalid", "must be in upper case")
	// Hexadecimal validates if a string is a valid hexadecimal number
	Hexadecimal = validation.NewStringValidator(govalidator.IsHexadecimal, "validation_hexadecimal_invalid", "must be a valid hexadecimal number")
	// HexColor validates if a string is a valid hexadecimal color code
	HexColor = validation.NewStringValidator(govalidator.IsHexcolor, "validation_hex_color_invalid", "must be a valid hexadecimal color code")
	// RGBColor validates if a string is a valid RGB color in the form of rgb(R, G, B)
	RGBColor = validation.NewStringValidator(govalidator.IsRGBcolor, "validation_rgb_color_invalid", "must be a valid RGB color code")
	// Int validates if a string is a valid integer number
	Int = validation.NewStringValidator(govalidator.IsInt, "validation_int_invalid", "must be an integer number")
	// Float validates if a string is a floating point number
	Float = validation.NewStringValidator(govalidator.IsFloat, "validation_float_invalid", "must be a floating point number")
	// UUIDv3 validates if a string is a valid version 3 UUID
	UUIDv3 = validation.NewStringValidator(govalidator.IsUUIDv3, "validation_uuid_v3_invalid", "must be a valid UUID v3")
	// UUIDv4 validates if a string is a valid version 4 UUID
	UUIDv4 = validation.NewStringValidator(govalidator.IsUUIDv4, "validation_uuid_v4_invalid", "must be a valid UUID v4")
	// UUIDv5 validates if a string is a valid version 5 UUID
	UUIDv5 = validation.NewStringValidator(govalidator.IsUUIDv5, "validation_uuid_v5_invalid", "must be a valid UUID v5")
	// UUID validates if a string is a valid UUID
	UUID = validation.NewStringValidator(govalidator.IsUUID, "validation_uuid_invalid", "must be a valid UUID")
	// CreditCard validates if a string is a valid credit card number
	CreditCard = validation.NewStringValidator(govalidator.IsCreditCard, "validation_credit_card_invalid", "must be a valid credit card number")
	// ISBN10 validates if a string is an ISBN version 10
	ISBN10 = validation.NewStringValidator(govalidator.IsISBN10, "validation_isbn_10_invalid", "must be a valid ISBN-10")
	// ISBN13 validates if a string is an ISBN version 13
	ISBN13 = validation.NewStringValidator(govalidator.IsISBN13, "validation_isbn_13_invalid", "must be a valid ISBN-13")
	// ISBN validates if a string is an ISBN (either version 10 or 13)
	ISBN = validation.NewStringValidator(isISBN, "validation_isbn_invalid", "must be a valid ISBN")
	// JSON validates if a string is in valid JSON format
	JSON = validation.NewStringValidator(govalidator.IsJSON, "validation_json_invalid", "must be in valid JSON format")
	// ASCII validates if a string contains ASCII characters only
	ASCII = validation.NewStringValidator(govalidator.IsASCII, "validation_ascii_invalid", "must contain ASCII characters only")
	// PrintableASCII validates if a string contains printable ASCII characters only
	PrintableASCII = validation.NewStringValidator(govalidator.IsPrintableASCII, "validation_printable_ascii_invalid", "must contain printable ASCII characters only")
	// Multibyte validates if a string contains multibyte characters
	Multibyte = validation.NewStringValidator(govalidator.IsMultibyte, "validation_multibyte_invalid", "must contain multibyte characters")
	// FullWidth validates if a string contains full-width characters
	FullWidth = validation.NewStringValidator(govalidator.IsFullWidth, "validation_full_width_invalid", "must contain full-width characters")
	// HalfWidth validates if a string contains half-width characters
	HalfWidth = validation.NewStringValidator(govalidator.IsHalfWidth, "validation_half_width_invalid", "must contain half-width characters")
	// VariableWidth validates if a string contains both full-width and half-width characters
	VariableWidth = validation.NewStringValidator(govalidator.IsVariableWidth, "validation_variable_width_invalid", "must contain both full-width and half-width characters")
	// Base64 validates if a string is encoded in Base64
	Base64 = validation.NewStringValidator(govalidator.IsBase64, "validation_base64_invalid", "must be encoded in Base64")
	// DataURI validates if a string is a valid base64-encoded data URI
	DataURI = validation.NewStringValidator(govalidator.IsDataURI, "validation_data_uri_invalid", "must be a Base64-encoded data URI")
	// E164 validates if a string is a valid ISO3166 Alpha 2 country code
	E164 = validation.NewStringValidator(isE164Number, "validation_e164_invalid", "must be a valid E164 number")
	// CountryCode2 validates if a string is a valid ISO3166 Alpha 2 country code
	CountryCode2 = validation.NewStringValidator(govalidator.IsISO3166Alpha2, "validation_country_code_2_invalid", "must be a valid two-letter country code")
	// CountryCode3 validates if a string is a valid ISO3166 Alpha 3 country code
	CountryCode3 = validation.NewStringValidator(govalidator.IsISO3166Alpha3, "validation_country_code_3_invalid", "must be a valid three-letter country code")
	// CurrencyCode validates if a string is a valid IsISO4217 currency code.
	CurrencyCode = validation.NewStringValidator(govalidator.IsISO4217, "validation_currency_code_invalid", "must be valid ISO 4217 currency code")
	// DialString validates if a string is a valid dial string that can be passed to Dial()
	DialString = validation.NewStringValidator(govalidator.IsDialString, "validation_dial_string_invalid", "must be a valid dial string")
	// MAC validates if a string is a MAC address
	MAC = validation.NewStringValidator(govalidator.IsMAC, "validation_mac_invalid", "must be a valid MAC address")
	// IP validates if a string is a valid IP address (either version 4 or 6)
	IP = validation.NewStringValidator(govalidator.IsIP, "validation_ip_invalid", "must be a valid IP address")
	// IPv4 validates if a string is a valid version 4 IP address
	IPv4 = validation.NewStringValidator(govalidator.IsIPv4, "validation_ipv4_invalid", "must be a valid IPv4 address")
	// IPv6 validates if a string is a valid version 6 IP address
	IPv6 = validation.NewStringValidator(govalidator.IsIPv6, "validation_ipv6_invalid", "must be a valid IPv6 address")
	// Subdomain validates if a string is valid subdomain
	Subdomain = validation.NewStringValidator(isSubdomain, "validation_sub_domain_invalid", "must be a valid subdomain")
	// Domain validates if a string is valid domain
	Domain = validation.NewStringValidator(isDomain, "validation_domain_invalid", "must be a valid domain")
	// DNSName validates if a string is valid DNS name
	DNSName = validation.NewStringValidator(govalidator.IsDNSName, "validation_dns_name_invalid", "must be a valid DNS name")
	// Host validates if a string is a valid IP (both v4 and v6) or a valid DNS name
	Host = validation.NewStringValidator(govalidator.IsHost, "validation_host_invalid", "must be a valid IP address or DNS name")
	// Port validates if a string is a valid port number
	Port = validation.NewStringValidator(govalidator.IsPort, "validation_port_invalid", "must be a valid port number")
	// MongoID validates if a string is a valid Mongo ID
	MongoID = validation.NewStringValidator(govalidator.IsMongoID, "validation_mongo_id_invalid", "must be a valid hex-encoded MongoDB ObjectId")
	// Latitude validates if a string is a valid latitude
	Latitude = validation.NewStringValidator(govalidator.IsLatitude, "validation_latitude_invalid", "must be a valid latitude")
	// Longitude validates if a string is a valid longitude
	Longitude = validation.NewStringValidator(govalidator.IsLongitude, "validation_longitude_invalid", "must be a valid longitude")
	// SSN validates if a string is a social security number (SSN)
	SSN = validation.NewStringValidator(govalidator.IsSSN, "validation_ssn_invalid", "must be a valid social security number")
	// Semver validates if a string is a valid semantic version
	Semver = validation.NewStringValidator(govalidator.IsSemver, "validation_semver_invalid", "must be a valid semantic version")
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
