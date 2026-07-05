package merge

import (
	"bufio"
)

// cmpLines returns the index of the lowest string in lines per cmp, skipping nil entries.
// Returns -1 if all entries are nil.
func cmpLines(lines []*string, cmp func(a, b string) int) int {
	lowest := -1
	for i, s := range lines {
		if s == nil {
			continue
		}
		if lowest == -1 || cmp(*s, *lines[lowest]) < 0 {
			lowest = i
		}
	}
	return lowest
}

// advanceStream reads the next line from each scanner at the given indexes, populating lines.
// Scanners that hit EOF are set to nil in the returned slice; lines entries for EOF scanners are set to nil.
func advanceStream(inputs []*bufio.Scanner, indexes []int, lines []*string) ([]*string, []*bufio.Scanner, error) {
	for _, idx := range indexes {
		inp := inputs[idx]
		if !inp.Scan() {
			if inp.Err() != nil {
				return lines, inputs, inp.Err()
			}
			lines[idx] = nil
			inputs[idx] = nil
		} else {
			s := inp.Text()
			lines[idx] = &s
		}
	}

	return lines, inputs, nil
}

// MergeStreams reads inputs a line at a time, and writes them to the out writer in low to high order as defined by the cmp function.
// Returns error or nil
// Each input must be already ordered, unique and not nil or the behavior is undefined
func MergeStreams(out *bufio.Writer, cmp func(a, b string) int, inputs ...*bufio.Scanner) error {

	lines := make([]*string, len(inputs))

	// Advance all files to read the first line of each
	idx := make([]int, len(inputs))
	for i := 0; i < len(inputs); i++ {
		idx[i] = i
	}
	var err error
	lines, inputs, err = advanceStream(inputs, idx, lines)
	if err != nil {
		return err
	}

	// Process until all empty
	for {
		// Next line
		advanceIdx := cmpLines(lines, cmp)

		// All strings nil then finished
		if advanceIdx == -1 {
			return nil
		}

		// Write the line
		if _, err := out.WriteString(*lines[advanceIdx] + "\n"); err != nil {
			return err
		}
		// Advance the stream that we wrote from
		lines, inputs, err = advanceStream(inputs, []int{advanceIdx}, lines)
		if err != nil {
			return err
		}
	}
}
