package diff

import (
	"bytes"
	"fmt"
)

type hunk struct {
	offsetA, offsetB int
	countA, countB   int
	lines            []Line
}

type Patch []Line

func (patch Patch) String() string {

	buf := bytes.NewBuffer(nil)
	//hunk := bytes.NewBuffer(nil)
	var h hunk

	beforeContext := true
	afterContext := 0
	lastAdded := -1

	for i, p := range patch {
		switch p.Action {
		case DiffDelete, DiffInsert:
			afterContext = 3

			if beforeContext {
				before := max(0, max(i-3, lastAdded))
				for j := before; j < i; j++ {
					if len(h.lines) == 0 { // start of hunk
						h.offsetA = patch[j].IndexLeft + 1
						h.offsetB = patch[j].IndexRight + 1
					}

					//fmt.Fprintf(hunk, "%s%s\n", patch[j].Action, patch[j].Content)
					h.lines = append(h.lines, patch[j])

					h.countA++
					h.countB++
				}
				beforeContext = false
			}

			//fmt.Fprintf(hunk, "%s%s\n", p.Action, p.Content)
			h.lines = append(h.lines, p)

			if p.Side == Left {
				h.countA++
			} else {
				h.countB++
			}
		default:
			if afterContext > 0 {
				//fmt.Fprintf(hunk, "%s%s\n", p.Action, p.Content)
				h.lines = append(h.lines, p)

				h.countA++
				h.countB++
				lastAdded = i
				afterContext--
				if afterContext == 0 { // hunk done
					fmt.Fprintf(buf, "@@ -%d,%d +%d,%d @@\n", h.offsetA, h.countA, h.offsetB, h.countB)
					for _, l := range h.lines {
						fmt.Fprintf(buf, "%s%s\n", l.Action, l.Content)
					}
					h.lines = make([]Line, 0, 50)
					h.countA, h.countB = 0, 0
					beforeContext = true
				}
			}
		}
	}

	fmt.Fprintf(buf, "@@ -%d,%d +%d,%d @@\n", h.offsetA, h.countA, h.offsetB, h.countB)
	for _, l := range h.lines {
		fmt.Fprintf(buf, "%s%s\n", l.Action, l.Content)
	}

	return buf.String()
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
