package main

import (
	"fmt"

	"github.com/antonmedv/expr/docgen"
	"github.com/antonmedv/expr/test/real_world"
)

func main() {
	doc := docgen.CreateDoc(real_world.NewEnv())

	fmt.Println(doc.Markdown())
}
