//go:build integration

package tests

import "testing"

func Hello() string {
	return "Hello, world"
}
func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
