package file

import (
	"testing"
)

// TestErrors reporting and recording.
func TestErrors(t *testing.T) {
	source := NewSource("a.b\n&&arg(missing, paren")
	errors := NewErrors(source)
	errors.ReportError(NewLocation(1, 1), "No such field")
	if len(errors.GetErrors()) != 1 {
		t.Errorf("%s first error not recorded", t.Name())
	}
	errors.ReportError(
		NewLocation(2, 20),
		"syntax error, missing paren",
	)
	if len(errors.GetErrors()) != 2 {
		t.Errorf("%s second error not recorded", t.Name())
	}
	got := errors.Error()
	want :=
		"No such field (1:2)\n" +
			" | a.b\n" +
			" | .^\n" +
			"syntax error, missing paren (2:21)\n" +
			" | &&arg(missing, paren\n" +
			" | ....................^"
	if got != want {
		t.Errorf("%s got %s, wanted %s", t.Name(), got, want)
	}
}

func TestErrors_WideAndNarrowCharacters(t *testing.T) {
	source := NewSource("你好吗\n我b很好\n")
	errors := NewErrors(source)
	errors.ReportError(NewLocation(2, 3), "Unexpected character '好'")

	got := errors.Error()
	want := "Unexpected character '好' (2:4)\n" +
		" | 我b很好\n" +
		" | ．.．＾"
	if got != want {
		t.Errorf("%s got %s, wanted %s", t.Name(), got, want)
	}
}

func TestErrors_WideAndNarrowCharacters_Emojis(t *testing.T) {
	source := NewSource("      '😁' in ['😁', '😑', '😦'] && in.😁")
	errors := NewErrors(source)
	errors.ReportError(NewLocation(1, 32), "syntax error: extraneous input 'in' expecting {'[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}")
	errors.ReportError(NewLocation(1, 35), "syntax error: token recognition error at: '😁'")
	errors.ReportError(NewLocation(1, 36), "syntax error: missing IDENTIFIER at '<EOF>'")
	got := errors.Error()
	want := "syntax error: extraneous input 'in' expecting {'[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER} (1:33)\n" +
		" |       '😁' in ['😁', '😑', '😦'] && in.😁\n" +
		" | .......．.......．....．....．......^\n" +
		"syntax error: token recognition error at: '😁' (1:36)\n" +
		" |       '😁' in ['😁', '😑', '😦'] && in.😁\n" +
		" | .......．.......．....．....．.........＾\n" +
		"syntax error: missing IDENTIFIER at '<EOF>' (1:37)\n" +
		" |       '😁' in ['😁', '😑', '😦'] && in.😁\n" +
		" | .......．.......．....．....．.........．^"
	if got != want {
		t.Errorf("%s got:\n%s\n\nwanted:\n%s", t.Name(), got, want)
	}
}
