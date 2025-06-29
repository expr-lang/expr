package parser

import "testing"

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
	p := new(Parser)
	for i := 0; i < b.N; i++ {
		p.Parse(source, nil)
	}
}
