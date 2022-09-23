package file

import (
	"testing"
)

const (
	unexpectedSnippet = "%s got snippet '%s', want '%v'"
	snippetNotFound   = "%s snippet not found, wanted '%v'"
	snippetFound      = "%s snippet found at Line %d, wanted none"
)

func TestStringSource_SnippetMultiLine(t *testing.T) {
	source := NewSource("hello\nworld\nmy\nbub\n")
	if str, found := source.Snippet(1); !found {
		t.Errorf(snippetNotFound, t.Name(), 1)
	} else if str != "hello" {
		t.Errorf(unexpectedSnippet, t.Name(), str, "hello")
	}
	if str2, found := source.Snippet(2); !found {
		t.Errorf(snippetNotFound, t.Name(), 2)
	} else if str2 != "world" {
		t.Errorf(unexpectedSnippet, t.Name(), str2, "world")
	}
	if str3, found := source.Snippet(3); !found {
		t.Errorf(snippetNotFound, t.Name(), 3)
	} else if str3 != "my" {
		t.Errorf(unexpectedSnippet, t.Name(), str3, "my")
	}
	if str4, found := source.Snippet(4); !found {
		t.Errorf(snippetNotFound, t.Name(), 4)
	} else if str4 != "bub" {
		t.Errorf(unexpectedSnippet, t.Name(), str4, "bub")
	}
	if str5, found := source.Snippet(5); !found {
		t.Errorf(snippetNotFound, t.Name(), 5)
	} else if str5 != "" {
		t.Errorf(unexpectedSnippet, t.Name(), str5, "")
	}
}

func TestStringSource_SnippetSingleLine(t *testing.T) {
	source := NewSource("hello, world")
	if str, found := source.Snippet(1); !found {
		t.Errorf(snippetNotFound, t.Name(), 1)

	} else if str != "hello, world" {
		t.Errorf(unexpectedSnippet, t.Name(), str, "hello, world")
	}
	if str2, found := source.Snippet(2); found {
		t.Errorf(snippetFound, t.Name(), 2)
	} else if str2 != "" {
		t.Errorf(unexpectedSnippet, t.Name(), str2, "")
	}
}
