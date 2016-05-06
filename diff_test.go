package diff

import (
	"fmt"
	"strings"
	"testing"
)

func callDiff(aText, bText, sep string) {
	a := strings.Split(aText, sep)
	b := strings.Split(bText, sep)

	//fmt.Println(aText)
	//fmt.Println(bText)

	result := calc(stringInterface{a, b})

	for _, r := range result {
		if r.Side == Right {
			fmt.Println(r.Action, r.Side, b[r.IndexRight])
		} else {
			fmt.Println(r.Action, r.Side, a[r.IndexLeft])
		}
	}
}

func callDiffF(aText, bText, sep string) {
	a := strings.Split(aText, sep)
	b := strings.Split(bText, sep)

	//fmt.Println(aText)
	//fmt.Println(bText)

	result := calc(stringInterface{a, b})

	fmt.Println(Patch{a, b, result})
}

func TestQuality(t *testing.T) {
	a := `
void func1() {
	x += 1
}

void func2() {
	x += 2
}`
	b := `
void func1() {
	x += 1
}

void functhreehalves() {
	x += 1.5
}

void func2() {
	x += 2
}`

	callDiff(a, b, "\n")
}

func TestMidInsert(t *testing.T) {
	a := "a z b x y w"
	b := "a z b c d w b x y w"

	callDiff(a, b, " ")
}

func TestEndChange(t *testing.T) {
	a := "a z b x y w"
	b := "a z b x y c"

	callDiff(a, b, " ")
}

func TestEndChange2(t *testing.T) {
	a := "a z b x a w"
	b := "a z b x b c"

	callDiff(a, b, " ")
}

func TestEndChange3(t *testing.T) {
	a := "a z b x a w n t"
	b := "a z b x b c n t"

	callDiff(a, b, " ")
}

func TestEndChange4(t *testing.T) {
	a := "a z b x a w n"
	b := "a z b x b c n"

	callDiff(a, b, " ")
}

func TestEndChange5(t *testing.T) {
	a := "1 2 a z b x a w n"
	b := "a z b x b c n"

	callDiff(a, b, " ")
}

func TestEndChange6(t *testing.T) {
	a := "a z b x a w n"
	b := "1 2 a z b x b c n"

	callDiff(a, b, " ")
}

func TestWorst(t *testing.T) {
	a := "a b c d"
	b := "w x y z"

	callDiff(a, b, " ")
}

func TestEqual(t *testing.T) {
	a := "a b c d"
	b := "a b c d"

	callDiff(a, b, " ")
}

func TestBasic(t *testing.T) {
	a := "1 2 1"
	b := "1 3 1"

	callDiff(a, b, " ")
}

func TestEndEqual(t *testing.T) {
	a := "a b c d f g h j q z w"
	b := "a b c d e f g i j k r x y z w"

	callDiff(a, b, " ")
}

func TestWiki(t *testing.T) {
	atext := `This part of the
document has stayed the
same from version to
version.  It shouldn't
be shown if it doesn't
change.  Otherwise, that
would not be helping to
compress the size of the
changes.

This paragraph contains
text that is outdated.
It will be deleted in the
near future.

It is important to spell
check this dokument. On
the other hand, a
misspelled word isn't
the end of the world.
Nothing in the rest of
this paragraph needs to
be changed. Things can
be added after it.`

	btext := `This is an important
notice! It should
therefore be located at
the beginning of this
document!

This part of the
document has stayed the
same from version to
version.  It shouldn't
be shown if it doesn't
change.  Otherwise, that
would not be helping to
compress anything.

It is important to spell
check this document. On
the other hand, a
misspelled word isn't
the end of the world.
Nothing in the rest of
this paragraph needs to
be changed. Things can
be added after it.

This paragraph contains
important new additions
to this document.`

	callDiffF(atext, btext, "\n")

}

func TestFiles(t *testing.T) {
	//	fmt.Println(Files("a.txt", "b.txt"))
}
