// Formats a table.
//
// a | | c |d
//
// |bbb| |ddd
//
//
// a |     | c |   d
//
//   | bbb |   | ddd
//
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	var lines []string
	var cols []int
	for s, e := r.ReadString('\n'); e == nil; s, e = r.ReadString('\n') {
		lines = append(lines, s)

		if len(strings.TrimSpace(s)) == 0 {
			continue
		}

		parts := strings.Split(s, "|")
		for i, p := range parts {
			if i == len(cols) {
				cols = append(cols, 0)
			}
			x := len(strings.TrimSpace(p))
			if x > cols[i] {
				cols[i] = x
			}

		}
	}

	for _, s := range lines {
		if len(strings.TrimSpace(s)) == 0 {
			fmt.Println()
			continue
		}

		parts := strings.Split(s, "|")
		for i, p := range parts {
			if i > 0 {
				fmt.Printf(" | ")
			}

			x := strings.TrimSpace(p)
			for j := len(x); j < cols[i]; j += 1 {
				fmt.Printf(" ")
			}
			fmt.Printf(x)
		}
		fmt.Println()
	}
}
