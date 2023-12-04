package main

import (
	"fmt"

	"github.com/expr-lang/expr/docgen"
	"github.com/expr-lang/expr/test/real_world"
)

func main() {
	doc := docgen.CreateDoc(real_world.NewEnv())

	fmt.Println(doc.Markdown())
}
