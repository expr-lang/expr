package main

import (
	"strings"
)

// CodeBlock holds the optional title and content of a code block.
type CodeBlock struct {
	Title   string
	Content string
}

func extractCodeBlocksWithTitle(markdown string) []CodeBlock {
	var blocks []CodeBlock
	var currentBlock []string
	var currentTitle string
	inBlock := false

	// Split the markdown into lines.
	lines := strings.Split(markdown, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check if the line starts with a code block fence.
		if strings.HasPrefix(trimmed, "```") {
			// If already inside a code block, this marks its end.
			if inBlock {
				blocks = append(blocks, CodeBlock{
					Title:   currentTitle,
					Content: strings.Join(currentBlock, "\n"),
				})
				currentBlock = nil
				inBlock = false
				currentTitle = ""
			} else {
				// Not in a block: starting a new code block.
				// Look backwards for the closest non-empty line that is not a code fence.
				title := ""
				for j := i - 1; j >= 0; j-- {
					prev := strings.TrimSpace(lines[j])
					if prev == "" || strings.HasPrefix(prev, "```") {
						continue
					}
					title = prev
					break
				}
				currentTitle = title
				inBlock = true
			}
			// Skip the fence line.
			continue
		}

		// If inside a code block, add the line.
		if inBlock {
			currentBlock = append(currentBlock, line)
		}
	}
	return blocks
}
