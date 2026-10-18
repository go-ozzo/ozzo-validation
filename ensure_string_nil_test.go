package validation

import "testing"

func TestEnsureStringNilReturnsError(t *testing.T) {
	if _, err := EnsureString(nil); err == nil {
		t.Fatal("EnsureString(nil) must report an unsupported input type")
	}
	if got, err := EnsureString([]byte(nil)); got != "" || err != nil {
		t.Fatalf("nil byte slice = %q, %v; want empty string, nil", got, err)
	}
}
