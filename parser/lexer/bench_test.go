package lexer

import (
	"testing"

	"github.com/expr-lang/expr/file"
)

func BenchmarkParser(b *testing.B) {
	const source = `
		/*
			Showing worst case scenario
		*/
		let value = trim("contains escapes \n\"\\ \U0001F600 and non ASCII ñ"); // inline comment
		len(value) == 0x2A
		// let's introduce an error too
		whatever
	`
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Lex(file.NewSource(source))
	}
}
