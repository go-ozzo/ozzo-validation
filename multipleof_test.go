// Copyright 2016 Qiang Xue, Google LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleof(t *testing.T) {
	r := MultipleOf(10)
	assert.Equal(t, "must be multiple of 10", r.Validate(11).Error())
	assert.Equal(t, nil, r.Validate(20))
}

func Test_Multipleof_Error(t *testing.T) {
	r := MultipleOf(10)
	assert.Equal(t, "must be multiple of 10", r.message)
}
