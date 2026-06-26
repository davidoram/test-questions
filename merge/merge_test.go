package merge

import (
	"bufio"
	"fmt"
	"strings"
)

func ExampleMergeStreams() {
	a := bufio.NewScanner(strings.NewReader("apple\ncherry\n"))
	b := bufio.NewScanner(strings.NewReader("banana\ndate\n"))

	var out strings.Builder
	w := bufio.NewWriter(&out)

	_ = MergeStreams(w, strings.Compare, a, b)
	w.Flush()

	fmt.Print(out.String())
	// Output:
	// apple
	// banana
	// cherry
	// date
}

func ExampleMergeStreams_nil() {
	a := bufio.NewScanner(strings.NewReader("apple\ncherry\n"))
	b := bufio.NewScanner(strings.NewReader(""))
	c := bufio.NewScanner(strings.NewReader("banana"))

	var out strings.Builder
	w := bufio.NewWriter(&out)

	_ = MergeStreams(w, strings.Compare, a, b, c)
	w.Flush()

	fmt.Print(out.String())
	// Output:
	// apple
	// banana
	// cherry
}
