package handler

import "testing"

func TestValidateInput(t *testing.T) {
	if !validateInput("Hello\nWorld") {
		t.Fatal("expected valid input")
	}

	if validateInput("Hello🙂") {
		t.Fatal("expected invalid input")
	}

	if !validateInput("Hello\\nWorld") {
		t.Fatal("expected backslash input to be allowed")
	}
}

func TestInterpretEscapeSequences(t *testing.T) {
	got := interpretEscapeSequences("Hello\\nWorld")
	want := "Hello\nWorld"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
