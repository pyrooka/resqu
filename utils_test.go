package main

import "testing"

func TestStringInSlice(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		slice := []string{"foo", "bar"}
		found := stringInSlice(slice, "bar")
		if found != true {
			t.Error("expected true")
		}
	})

	t.Run("not found", func(t *testing.T) {
		slice := []string{"foo", "bar"}
		found := stringInSlice(slice, "foobar")
		if found != false {
			t.Error("expected false")
		}
	})
}
