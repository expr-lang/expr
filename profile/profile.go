package profile

import (
	"bytes"
	"fmt"
	"os"

	"github.com/google/pprof/profile"

	"github.com/expr-lang/expr/vm"
)

// GeneratePprofProfile generates a pprof-formatted profile file based on the Span structure.
// Parameters:
//   - rootSpan: The root Span structure, containing information about the expression's runtime and hierarchical structure.
//   - filePath: The path to save the generated pprof file.
//
// Returns:
//   - An error if an error occurs during the generation process; otherwise, returns nil.
func GeneratePprofProfile(rootSpan *vm.Span, filePath string) error {
	// Create a new pprof configuration file
	p := &profile.Profile{
		// Define the type and unit of the sampling period, here it's CPU time in nanoseconds
		PeriodType: &profile.ValueType{
			Type: "cpu",
			Unit: "nanoseconds",
		},
		// Define the type and unit of the sampling data, also CPU time in nanoseconds
		SampleType: []*profile.ValueType{
			{
				Type: "cpu",
				Unit: "nanoseconds",
			},
		},
	}

	// Create maps to store function and location information
	// The key is the function name, and the value is the profile.Function structure
	functionMap := make(map[string]*profile.Function)
	// The key is a combination of the function name and the expression, and the value is the profile.Location structure
	locationMap := make(map[string]*profile.Location)

	// Recursively traverse the Span structure
	var traverse func(span *vm.Span, parentLocation *profile.Location)
	traverse = func(span *vm.Span, parentLocation *profile.Location) {
		// Get or create function and location information
		// Retrieve the function name from the Span structure
		functionName := span.Name
		// Check if the function name already exists in the functionMap
		if _, ok := functionMap[functionName]; !ok {
			// If not, create a new profile.Function structure
			functionMap[functionName] = &profile.Function{
				// Assign a unique ID
				ID: uint64(len(p.Function) + 1),
				// Function name
				Name: functionName,
				// System function name, same as the function name here
				SystemName: functionName,
			}
			// Add the new function information to the pprof configuration file
			p.Function = append(p.Function, functionMap[functionName])
		}

		// Generate the key for the location information, a combination of the function name and the expression
		locationKey := fmt.Sprintf("%s:%s", functionName, span.Expression)
		// Check if the location information already exists in the locationMap
		if _, ok := locationMap[locationKey]; !ok {
			// If not, create a new profile.Location structure
			locationMap[locationKey] = &profile.Location{
				// Assign a unique ID
				ID: uint64(len(p.Location) + 1),
				// Line number information for the location
				Line: []profile.Line{
					{
						// Associated function information
						Function: functionMap[functionName],
						// Line number is 1
						Line: 1,
					},
				},
			}
			// Add the new location information to the pprof configuration file
			p.Location = append(p.Location, locationMap[locationKey])
		}

		// Create sample information
		sample := &profile.Sample{
			// The value of the sample, i.e., the duration of the Span converted to nanoseconds
			Value: []int64{int64(span.Duration)},
			// Location information for the sample
			Location: []*profile.Location{
				locationMap[locationKey],
			},
		}
		// If there is parent location information, add it to the sample's location information
		if parentLocation != nil {
			sample.Location = append([]*profile.Location{parentLocation}, sample.Location...)
		}
		// Add the new sample information to the pprof configuration file
		p.Sample = append(p.Sample, sample)

		// Recursively process child Spans
		for _, child := range span.Children {
			// Recursively call the traverse function to process child Spans
			traverse(child, locationMap[locationKey])
		}
	}

	// Start traversing from the root Span
	traverse(rootSpan, nil)

	// Write to the pprof file
	// Create a byte buffer to store the content of the pprof file
	var buf bytes.Buffer
	// Write the pprof configuration file to the buffer
	if err := p.Write(&buf); err != nil {
		// If an error occurs during the writing process, return the error information
		return err
	}

	// Write the content of the buffer to the specified file
	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		// If an error occurs during the file writing process, return the error information
		return err
	}

	// No errors occurred, return nil
	return nil
}
