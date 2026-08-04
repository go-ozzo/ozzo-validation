package validation

import (
	"reflect"
	"sync"
)

type fieldCacheEntry struct {
	field  reflect.StructField
	offset uintptr
}

var (
	fieldCacheMu sync.RWMutex
	fieldCache   = make(map[reflect.Type][]fieldCacheEntry)
)

func getFieldEntries(structType reflect.Type) []fieldCacheEntry {
	fieldCacheMu.RLock()
	entries, ok := fieldCache[structType]
	fieldCacheMu.RUnlock()
	if ok {
		return entries
	}

	fieldCacheMu.Lock()
	defer fieldCacheMu.Unlock()

	if entries, ok = fieldCache[structType]; ok {
		return entries
	}

	entries = buildFieldEntries(structType, 0)
	fieldCache[structType] = entries
	return entries
}

func buildFieldEntries(structType reflect.Type, baseOffset uintptr) []fieldCacheEntry {
	var entries []fieldCacheEntry
	for i := 0; i < structType.NumField(); i++ {
		sf := structType.Field(i)
		entries = append(entries, fieldCacheEntry{
			field:  sf,
			offset: baseOffset + sf.Offset,
		})
		if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
			entries = append(entries, buildFieldEntries(sf.Type, baseOffset+sf.Offset)...)
		}
	}
	return entries
}

func findStructFieldCached(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	fieldPtr := fieldValue.Pointer()
	structPtr := structValue.UnsafeAddr()
	fieldOffset := fieldPtr - structPtr

	entries := getFieldEntries(structValue.Type())

	for i := range entries {
		if entries[i].offset == fieldOffset && entries[i].field.Type == fieldValue.Elem().Type() {
			return &entries[i].field
		}
	}

	return nil
}
