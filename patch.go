package diff

import (
	"bytes"
	"fmt"
)

type Patch []Line

func (patch Patch) String() string {

	buf := bytes.NewBuffer(nil)
	hunk := bytes.NewBuffer(nil)

	beforeContext := true
	afterContext := 0
	lastAdded := -1
	aStart, aCnt, bStart, bCnt := 1, 0, 1, 0

	for i, p := range patch {
		switch p.Action {
		case DiffDelete, DiffInsert:
			afterContext = 3

			if beforeContext {
				before := max(0, max(i-3, lastAdded))
				for j := before; j < i; j++ {
					if hunk.Len() == 0 { // start of hunk
						aStart = patch[j].IndexLeft + 1
						bStart = patch[j].IndexRight + 1
					}

					fmt.Fprintf(hunk, "%s%s\n", patch[j].Action, patch[j].Content)
					aCnt++
					bCnt++
				}
				beforeContext = false
			}

			fmt.Fprintf(hunk, "%s%s\n", p.Action, p.Content)

			if p.Side == Left {
				aCnt++
			} else {
				bCnt++
			}
		default:

			if afterContext > 0 {
				fmt.Fprintf(hunk, "%s%s\n", p.Action, p.Content)
				aCnt++
				bCnt++
				lastAdded = i
				afterContext--
				if afterContext == 0 { // hunk done
					fmt.Fprintf(buf, "@@ -%d,%d +%d,%d @@\n", aStart, aCnt, bStart, bCnt)
					buf.Write(hunk.Bytes())
					hunk.Reset()
					aCnt, bCnt = 0, 0
					beforeContext = true
				}
			}

		}
	}

	fmt.Fprintf(buf, "@@ -%d,%d +%d,%d @@\n", aStart, aCnt, bStart, bCnt)
	buf.Write(hunk.Bytes())

	return buf.String()
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
