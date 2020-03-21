# DocGen

This package provides documentation generator with JSON or Markdown output.

## Usage

Create a file and put next code into it. 

```go
package main

import (
	"encoding/json"
	"fmt"
  
	"github.com/antonmedv/expr/docgen"
)

func main() {
	// TODO: Replace env with your own types.
	doc := docgen.CreateDoc(env)
  
	buf, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
```

Run `go run your_file.go`. Documentation will be printed in JSON format.

## Markdown

To generate markdown documentation: 

```go
package main

import "github.com/antonmedv/expr/docgen"

func main() {
	// TODO: Replace env with your own types.
	doc := docgen.CreateDoc(env)

	print(doc.Markdown())
}
```
