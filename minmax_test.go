// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	date0 := time.Time{}
	date20000101 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	date20001201 := time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC)
	date20000601 := time.Date(2000, 6, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		tag       string
		threshold interface{}
		exclusive bool
		value     interface{}
		err       string
	}{
		// int cases
		{"t1.1", 1, false, 1, ""},
		{"t1.2", 1, false, 2, ""},
		{"t1.3", 1, false, -1, "must be no less than 1"},
		{"t1.4", 1, false, 0, ""},
		{"t1.5", 1, true, 1, "must be greater than 1"},
		{"t1.6", 1, false, "1", "cannot convert string to int64"},
		{"t1.7", "1", false, 1, "type not supported: string"},
		// uint cases
		{"t2.1", uint(2), false, uint(2), ""},
		{"t2.2", uint(2), false, uint(3), ""},
		{"t2.3", uint(2), false, uint(1), "must be no less than 2"},
		{"t2.4", uint(2), false, uint(0), ""},
		{"t2.5", uint(2), true, uint(2), "must be greater than 2"},
		{"t2.6", uint(2), false, "1", "cannot convert string to uint64"},
		// float cases
		{"t3.1", float64(2), false, float64(2), ""},
		{"t3.2", float64(2), false, float64(3), ""},
		{"t3.3", float64(2), false, float64(1), "must be no less than 2"},
		{"t3.4", float64(2), false, float64(0), ""},
		{"t3.5", float64(2), true, float64(2), "must be greater than 2"},
		{"t3.6", float64(2), false, "1", "cannot convert string to float64"},
		// Time cases
		{"t4.1", date20000601, false, date20000601, ""},
		{"t4.2", date20000601, false, date20001201, ""},
		{"t4.3", date20000601, false, date20000101, "must be no less than 2000-06-01 00:00:00 +0000 UTC"},
		{"t4.4", date20000601, false, date0, ""},
		{"t4.5", date20000601, true, date20000601, "must be greater than 2000-06-01 00:00:00 +0000 UTC"},
		{"t4.6", date20000601, true, 1, "cannot convert int to time.Time"},
		{"t4.7", struct{}{}, false, 1, "type not supported: struct {}"},
		{"t4.8", date0, false, date20000601, ""},
		// Zero values are empty and skipped, as in v4.3.0 (#234); ValidateEmpty checks them
		{"t5.1", 0, false, 0, ""},
		{"t5.2", 0, true, 0, ""},
		{"t5.3", float64(0), true, float64(0), ""},
		{"t5.4", 1, false, 0, ""},
		// nil value tests
		{"t6.1", 1, false, nil, ""},
	}

	for _, test := range tests {
		r := Min(test.threshold)
		if test.exclusive {
			r = r.Exclusive()
		}
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestMinError(t *testing.T) {
	r := Min(10)
	assert.Equal(t, "must be no less than 10", r.Validate(9).Error())

	r = r.Error("123")
	assert.Equal(t, "123", r.err.Message())
}

func TestMax(t *testing.T) {
	date0 := time.Time{}
	date20000101 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	date20001201 := time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC)
	date20000601 := time.Date(2000, 6, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		tag       string
		threshold interface{}
		exclusive bool
		value     interface{}
		err       string
	}{
		// int cases
		{"t1.1", 2, false, 2, ""},
		{"t1.2", 2, false, 1, ""},
		{"t1.3", 2, false, 3, "must be no greater than 2"},
		{"t1.4", 2, false, 0, ""},
		{"t1.5", 2, true, 2, "must be less than 2"},
		{"t1.6", 2, false, "1", "cannot convert string to int64"},
		{"t1.7", "1", false, 1, "type not supported: string"},
		// uint cases
		{"t2.1", uint(2), false, uint(2), ""},
		{"t2.2", uint(2), false, uint(1), ""},
		{"t2.3", uint(2), false, uint(3), "must be no greater than 2"},
		{"t2.4", uint(2), false, uint(0), ""},
		{"t2.5", uint(2), true, uint(2), "must be less than 2"},
		{"t2.6", uint(2), false, "1", "cannot convert string to uint64"},
		// float cases
		{"t3.1", float64(2), false, float64(2), ""},
		{"t3.2", float64(2), false, float64(1), ""},
		{"t3.3", float64(2), false, float64(3), "must be no greater than 2"},
		{"t3.4", float64(2), false, float64(0), ""},
		{"t3.5", float64(2), true, float64(2), "must be less than 2"},
		{"t3.6", float64(2), false, "1", "cannot convert string to float64"},
		// Time cases
		{"t4.1", date20000601, false, date20000601, ""},
		{"t4.2", date20000601, false, date20000101, ""},
		{"t4.3", date20000601, false, date20001201, "must be no greater than 2000-06-01 00:00:00 +0000 UTC"},
		{"t4.4", date20000601, false, date0, ""},
		{"t4.5", date20000601, true, date20000601, "must be less than 2000-06-01 00:00:00 +0000 UTC"},
		{"t4.6", date20000601, true, 1, "cannot convert int to time.Time"},
		// nil value tests
		{"t5.1", 2, false, nil, ""},
	}

	for _, test := range tests {
		r := Max(test.threshold)
		if test.exclusive {
			r = r.Exclusive()
		}
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestMaxError(t *testing.T) {
	r := Max(10)
	assert.Equal(t, "must be no greater than 10", r.Validate(11).Error())

	r = r.Error("123")
	assert.Equal(t, "123", r.err.Message())
}

func TestThresholdRule_ErrorObject(t *testing.T) {
	r := Max(10)
	err := NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.err)
	assert.Equal(t, err.Code(), r.err.Code())
	assert.Equal(t, err.Message(), r.err.Message())
}

func TestThresholdRule_ValidateEmpty(t *testing.T) {
	date0 := time.Time{}
	date1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	zero := 0

	tests := []struct {
		tag   string
		rule  ThresholdRule
		value interface{}
		err   string
	}{
		{"int 0 below Min", Min(1).ValidateEmpty(), 0, "must be no less than 1"},
		{"uint 0 below Min", Min(uint(2)).ValidateEmpty(), uint(0), "must be no less than 2"},
		{"float 0 below Min", Min(float64(2)).ValidateEmpty(), float64(0), "must be no less than 2"},
		{"int 0 not above exclusive Min (#165)", Min(0).Exclusive().ValidateEmpty(), 0, "must be greater than 0"},
		{"float 0 not above exclusive Min", Min(float64(0)).Exclusive().ValidateEmpty(), float64(0), "must be greater than 0"},
		{"zero time below Min (#180)", Min(date1900).ValidateEmpty(), date0, "must be no less than 1900-01-01 00:00:00 +0000 UTC"},
		{"int 0 above Max", Max(-1).ValidateEmpty(), 0, "must be no greater than -1"},
		{"pointer to 0", Min(1).ValidateEmpty(), &zero, "must be no less than 1"},
		{"int 0 within range", Min(0).ValidateEmpty(), 0, ""},
		{"zero time within range", Max(date1900).ValidateEmpty(), date0, ""},
		{"nil pointer stays valid", Min(1).ValidateEmpty(), (*int)(nil), ""},
		{"empty value of the wrong type", Min(10).ValidateEmpty(), "", "cannot convert string to int64"},
		{"ValidateEmpty before Exclusive", Min(0).ValidateEmpty().Exclusive(), 0, "must be greater than 0"},
	}

	for _, test := range tests {
		assertError(t, test.err, test.rule.Validate(test.value), test.tag)
	}

	// Without ValidateEmpty, empty values are skipped.
	assert.NoError(t, Min(1).Validate(0))
	assert.NoError(t, Min(date1900).Validate(date0))
}
