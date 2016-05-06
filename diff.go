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

// Interface provides interface used in this package for comparing two collections
type Interface interface {
	Compare(a, b int) int
	LenA() int
	LenB() int
}

// Line represents change
type Line struct {
	IndexLeft  int
	IndexRight int
	Action     string
	Side       string
}

// String splits both sides of string and creates diff result comparing both side lines
func String(left, right, sep string) (result []Line) {
	a := strings.Split(left, sep)
	b := strings.Split(right, sep)

	return calc(stringInterface{a, b})
}

// calc returns all changes lines, from to change collection
func calc(input Interface) (result []Line) {
	i, j := 0, 0

	for i < input.LenA() && j < input.LenB() {

		if input.Compare(i, j) == 0 {
			result = append(result, Line{IndexLeft: i, IndexRight: j, Action: DiffSame})
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
			if x >= input.LenA() {
				y += x + 1 - input.LenA()
				x = input.LenA() - 1
			}

			count := 0
			for x >= i && y < input.LenB() {
				count++
				if input.Compare(x, y) == 0 {
					found = true

					// while searching for next matching line, we choose matching pair farther away from last line
					m := min(input.LenA()-x, input.LenB()-y)
					if m >= d {
						atx, aty = x, y
						d = m
					}
				}
				x--
				y++
			}

			if count == 0 {
				for i < input.LenA() {
					result = append(result, Line{Side: Left, Action: DiffDelete, IndexLeft: i})
					i++
				}

				for j < input.LenB() {
					result = append(result, Line{Side: Right, Action: DiffInsert, IndexRight: j})
					j++
				}
				break
			}

			if found || count == 1 {

				if !found && count == 1 {
					x = input.LenA()
					y = input.LenB()
				}

				if found {
					x, y = atx, aty
				}

				for k := i; k < x; k++ { // removed lines from left
					result = append(result, Line{Side: Left, Action: DiffDelete, IndexLeft: k})
				}
				for k := j; k < y; k++ { // added lines to right
					result = append(result, Line{Side: Right, Action: DiffInsert, IndexRight: k})
				}

				if found { // add line that was found same in both lists
					result = append(result, Line{IndexLeft: x, IndexRight: y, Action: DiffSame})
				}

				i = x + 1
				j = y + 1
				break
			}
		}
	}

	for k := i; k < input.LenA(); k++ {
		result = append(result, Line{Side: Left, Action: DiffDelete, IndexLeft: k})
	}

	for k := j; k < input.LenB(); k++ {
		result = append(result, Line{Side: Right, Action: DiffInsert, IndexRight: k})
	}

	return
}

// Files returns difference between two files
func Files(pathA, pathB string) (patch Patch, err error) {
	left, err := ioutil.ReadFile(pathA)
	if err != nil {
		return
	}

	right, err := ioutil.ReadFile(pathB)
	if err != nil {
		return
	}

	a := strings.Split(string(left), "\n")
	b := strings.Split(string(right), "\n")

	patch = Patch{a, b, calc(stringInterface{a, b})}
	return
}
