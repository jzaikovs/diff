package diff

import (
	"io/ioutil"
	"strings"
)

const (
	DiffSame   = " "
	DiffInsert = "+"
	DiffUpdate = "~"
	DiffDelete = "-"
)

const (
	Left  = "A"
	Right = "B"
)

// Line represents change
type Line struct {
	IndexLeft  int
	IndexRight int
	Action     string
	Content    string
	Side       string
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// String splits both sides of string and creates diff result comparing both side lines
func String(left, right, sep string) (result []Line) {
	a := strings.Split(left, sep)
	b := strings.Split(right, sep)

	return calc(a, b)
}

// calc returns all changes lines, from to change collection
func calc(linesA, linesB []string) (result []Line) {
	i, j := 0, 0

	for i < len(linesA) && j < len(linesB) {
		if linesA[i] == linesB[j] {
			result = append(result, Line{IndexLeft: i, IndexRight: j, Action: DiffSame, Content: linesA[i]})
			i++
			j++
			continue
		}

		found := false
		atx, aty, n, d := 0, 0, 0, 0

		for {
			n++        // each cycle will increase search radius
			y := j     // sarch will start at mismach line
			x := i + n // top-right corner of search

			// move down if x is out of old line count
			if x >= len(linesA) {
				y += x + 1 - len(linesA)
				x = len(linesA) - 1
			}

			count := 0
			for x >= i && y < len(linesB) {
				count++
				if linesA[x] == linesB[y] {
					found = true

					// while searching for next matching line, we choose matching pair farther away from last line
					m := min(len(linesA)-x, len(linesB)-y)
					if m >= d {
						atx, aty = x, y
						d = m
					}
				}
				x--
				y++
			}

			if count == 0 {
				for i < len(linesA) {
					result = append(result, Line{Side: Left, Action: DiffDelete, Content: linesA[i]})
					i++
				}

				for j < len(linesB) {
					result = append(result, Line{Side: Right, Action: DiffInsert, Content: linesB[j]})
					j++
				}
				break
			}

			if found || count == 1 {

				if !found && count == 1 {
					x = len(linesA)
					y = len(linesB)
				}

				if found {
					x, y = atx, aty
				}

				for k := i; k < x; k++ { // removed lines from left
					result = append(result, Line{Side: Left, Action: DiffDelete, Content: linesA[k]})
				}
				for k := j; k < y; k++ { // added lines to right
					result = append(result, Line{Side: Right, Action: DiffInsert, Content: linesB[k]})
				}

				if found { // add line that was found same in both lists
					result = append(result, Line{IndexLeft: x, IndexRight: y, Action: DiffSame, Content: linesA[x]})
				}

				i = x + 1
				j = y + 1
				break
			}
		}
	}

	for k := i; k < len(linesA); k++ {
		result = append(result, Line{Side: Left, Action: DiffDelete, Content: linesA[k]})
	}

	for k := j; k < len(linesB); k++ {
		result = append(result, Line{Side: Right, Action: DiffInsert, Content: linesB[k]})
	}

	return
}

// Files returns difference between two files
func Files(pathA, pathB string) (patch Patch, err error) {
	a, err := ioutil.ReadFile(pathA)
	if err != nil {
		return
	}

	b, err := ioutil.ReadFile(pathB)
	if err != nil {
		return
	}

	patch = String(string(a), string(b), "\n")
	return
}
