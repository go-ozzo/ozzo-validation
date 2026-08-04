package validation

import (
	"context"
	"regexp"
	"testing"
)

type benchStruct struct {
	Name   string
	Email  string
	Age    int
	Street string
	City   string
	Zip    string
}

var benchZipRegex = regexp.MustCompile(`^\d{5}$`)

func BenchmarkValidateStruct(b *testing.B) {
	s := benchStruct{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
		Street: "123 Main St",
		City:   "Springfield",
		Zip:    "12345",
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateStruct(&s,
			Field(&s.Name, Required, Length(1, 100)),
			Field(&s.Email, Required, Length(5, 255)),
			Field(&s.Age, Required, Min(1), Max(150)),
			Field(&s.Street, Required, Length(1, 200)),
			Field(&s.City, Required, Length(1, 100)),
			Field(&s.Zip, Required, Match(benchZipRegex)),
		)
	}
}

func BenchmarkValidateStruct_Invalid(b *testing.B) {
	s := benchStruct{}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateStruct(&s,
			Field(&s.Name, Required, Length(1, 100)),
			Field(&s.Email, Required, Length(5, 255)),
			Field(&s.Age, Required, Min(1), Max(150)),
			Field(&s.Street, Required, Length(1, 200)),
			Field(&s.City, Required, Length(1, 100)),
			Field(&s.Zip, Required, Length(5, 10)),
		)
	}
}

func BenchmarkValidateStructWithContext(b *testing.B) {
	ctx := context.Background()
	s := benchStruct{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
		Street: "123 Main St",
		City:   "Springfield",
		Zip:    "12345",
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateStructWithContext(ctx, &s,
			Field(&s.Name, Required, Length(1, 100)),
			Field(&s.Email, Required, Length(5, 255)),
			Field(&s.Age, Required, Min(1), Max(150)),
			Field(&s.Street, Required, Length(1, 200)),
			Field(&s.City, Required, Length(1, 100)),
			Field(&s.Zip, Required, Length(5, 10)),
		)
	}
}

func BenchmarkValidateMap(b *testing.B) {
	m := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Validate(m,
			Map(
				Key("name", Required, Length(1, 100)),
				Key("email", Required, Length(5, 255)),
				Key("age", Required, Min(1)),
			),
		)
	}
}

func BenchmarkValidateSimpleValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Validate("test@example.com", Required, Length(5, 255))
	}
}

func BenchmarkValidateWithConditional(b *testing.B) {
	s := benchStruct{
		Name:  "John",
		Email: "john@example.com",
		Age:   30,
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ValidateStruct(&s,
			Field(&s.Name, Required),
			Field(&s.Email, When(s.Name != "", Required, Length(5, 255))),
			Field(&s.Age, When(s.Age > 0, Min(1), Max(150))),
		)
	}
}

func BenchmarkEach_100Items(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = "value"
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Validate(items, Each(Required, Length(1, 50)))
	}
}

func BenchmarkEachUntilFirstError_100Items(b *testing.B) {
	items := make([]string, 100)
	items[0] = ""
	for i := 1; i < len(items); i++ {
		items[i] = "value"
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Validate(items, EachUntilFirstError(Required))
	}
}
