// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFieldEntries(t *testing.T) {
	type1 := reflect.TypeOf(Struct1{})
	entries := buildFieldEntries(type1, 0)

	byName := make(map[string]fieldCacheEntry)
	for _, e := range entries {
		byName[e.field.Name] = e
	}

	assert.Equal(t, type1.Field(0).Offset, byName["Field1"].offset)

	// fields promoted through a value-embedded anonymous struct are flattened
	struct2Field, ok := type1.FieldByName("Struct2")
	assert.True(t, ok)
	field21, ok := struct2Field.Type.FieldByName("Field21")
	assert.True(t, ok)
	assert.Contains(t, byName, "Field21")
	assert.Equal(t, struct2Field.Offset+field21.Offset, byName["Field21"].offset)
	assert.Contains(t, byName, "Field22")

	// fields promoted through a pointer-embedded anonymous struct are NOT
	// flattened; ValidateStruct relies on the findStructField fallback for
	// this case (see struct.go)
	type3 := reflect.TypeOf(Struct3{})
	entries3 := buildFieldEntries(type3, 0)
	byName3 := make(map[string]fieldCacheEntry)
	for _, e := range entries3 {
		byName3[e.field.Name] = e
	}
	assert.Contains(t, byName3, "Struct2")
	assert.Contains(t, byName3, "S1")
	assert.NotContains(t, byName3, "Field21")
}

func TestGetFieldEntries(t *testing.T) {
	typ := reflect.TypeOf(Struct1{})

	fieldCacheMu.Lock()
	delete(fieldCache, typ)
	fieldCacheMu.Unlock()

	entries := getFieldEntries(typ)
	assert.NotEmpty(t, entries)

	fieldCacheMu.RLock()
	cached, ok := fieldCache[typ]
	fieldCacheMu.RUnlock()
	assert.True(t, ok)
	assert.Equal(t, entries, cached)

	// a second call must hit the cache and return the same data
	again := getFieldEntries(typ)
	assert.Equal(t, entries, again)
}

func TestGetFieldEntries_Concurrent(t *testing.T) {
	type ConcurrentCacheStruct struct {
		A, B, C, D, E int
	}
	typ := reflect.TypeOf(ConcurrentCacheStruct{})

	fieldCacheMu.Lock()
	delete(fieldCache, typ)
	fieldCacheMu.Unlock()

	const goroutines = 50
	var wg sync.WaitGroup
	results := make([][]fieldCacheEntry, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = getFieldEntries(typ)
		}(i)
	}
	wg.Wait()

	for i := 1; i < goroutines; i++ {
		assert.Equal(t, results[0], results[i])
	}
}

func TestFindStructFieldCached(t *testing.T) {
	var s1 Struct1
	v1 := reflect.ValueOf(&s1).Elem()
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Field1)))
	assert.Nil(t, findStructFieldCached(v1, reflect.ValueOf(s1.Field2)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Field2)))
	assert.Nil(t, findStructFieldCached(v1, reflect.ValueOf(s1.Field3)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Field3)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Field4)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.field5)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Struct2)))
	assert.Nil(t, findStructFieldCached(v1, reflect.ValueOf(s1.S1)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.S1)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Field21)))
	assert.NotNil(t, findStructFieldCached(v1, reflect.ValueOf(&s1.Field22)))
	s2 := reflect.ValueOf(&s1.Struct2).Elem()
	assert.NotNil(t, findStructFieldCached(s2, reflect.ValueOf(&s1.Field21)))
	assert.NotNil(t, findStructFieldCached(s2, reflect.ValueOf(&s1.Field22)))

	// findStructFieldCached does not flatten fields promoted through a
	// pointer-embedded anonymous struct; ValidateStruct falls back to
	// findStructField to handle this case (see struct.go).
	s3 := Struct3{Struct2: &Struct2{}}
	v3 := reflect.ValueOf(&s3).Elem()
	assert.NotNil(t, findStructFieldCached(v3, reflect.ValueOf(&s3.Struct2)))
	assert.Nil(t, findStructFieldCached(v3, reflect.ValueOf(&s3.Field21)))
}
