package viemu

import (
	"testing"

	"github.com/iivvoo/novi/novi"
)

func AssertLinePos(t *testing.T, expLine, expPos, actLine, actPos int) {
	t.Helper()
	if expLine != actLine || expPos != actPos {
		t.Errorf("Expected (l,p) to be (%d, %d) but got (%d, %d)", expLine, expPos, actLine, actPos)
	}
}

func TestJumpWordForward(t *testing.T) {
	b := novi.BuildBuffer("This is the first line.", "", "  leading space",
		"trailing space   ", "last line")

	t.Run("Find first from start", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 0, 5, l, p)
	})
	t.Run("Find empty line", func(t *testing.T) {
		c := b.NewCursor(0, 19) // on 'i' in line
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 1, 0, l, p)
	})
	t.Run("Start from empty line", func(t *testing.T) {
		c := b.NewCursor(1, 0)
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 2, 2, l, p)
	})
	t.Run("Find from middle of word", func(t *testing.T) {
		c := b.NewCursor(0, 13) // on 'i' in first
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 0, 18, l, p)
	})
	t.Run("Find on next line", func(t *testing.T) {
		c := b.NewCursor(2, 11) // on 'a' in space
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 3, 0, l, p)
	})
	t.Run("Find end", func(t *testing.T) {
		c := b.NewCursor(4, 5) // on 'l' in line
		l, p := JumpWordForward(b, c)

		AssertLinePos(t, 4, 8, l, p)
	})
}

// Test 'B' behaviour
func TestJumpWordBackward(t *testing.T) {
	b := novi.BuildBuffer("This is the first line.", "", "  leading space",
		"trailing space   ", "last line")

	t.Run("Find first from end", func(t *testing.T) {
		c := b.NewCursor(4, 8)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 4, 5, l, p)
	})
	t.Run("Find empty line", func(t *testing.T) {
		c := b.NewCursor(2, 2) // on 'l' in leading
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 1, 0, l, p)
	})
	t.Run("Find from middle of word", func(t *testing.T) {
		c := b.NewCursor(0, 13) // on 'i' in first
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 12, l, p)
	})
	t.Run("Find on previous line", func(t *testing.T) {
		c := b.NewCursor(4, 0)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 3, 9, l, p) // 's' in space
	})
	t.Run("From/to first word on line", func(t *testing.T) {
		c := b.NewCursor(3, 3)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 3, 0, l, p)
	})
	t.Run("From/to first word on line", func(t *testing.T) {
		b := novi.BuildBuffer("First line", "    second line with spaces")
		c := b.NewCursor(1, 3) // at space character
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 6, l, p)
	})
	t.Run("Jump from empty", func(t *testing.T) {
		c := b.NewCursor(1, 0)
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 18, l, p)
	})
	t.Run("Jump previous line word continues", func(t *testing.T) {
		b := novi.BuildBuffer("foo bar", "4. bla 123")
		c := b.NewCursor(1, 0) // where we ended the previous test
		l, p := JumpWordBackward(b, c)

		AssertLinePos(t, 0, 4, l, p)
	})
}

// Tests "w" behaviour
func TestJumpForward(t *testing.T) {
	b := novi.BuildBuffer("This..isa line/;-with? separators", "", "  leading space",
		"https://github.com/some/repo.git?foo=a", "last line")

	t.Run("Find first from start", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 4, l, p)
	})
	t.Run("Find first from interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 4)
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 6, l, p)
	})
	t.Run("Find interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 10) // 'l' in line
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 14, l, p)
	})
	t.Run("Find first from interpunction, skipping whitespace", func(t *testing.T) {
		c := b.NewCursor(0, 21) // the ? after with
		l, p := JumpForward(b, c)

		AssertLinePos(t, 0, 23, l, p) // expect space to be skipped
	})
	t.Run("Find empty line", func(t *testing.T) {
		c := b.NewCursor(0, 23) // the 's' in separators
		l, p := JumpForward(b, c)

		AssertLinePos(t, 1, 0, l, p) // expect space to be skipped
	})
	t.Run("Find from empty line", func(t *testing.T) {
		c := b.NewCursor(1, 0) // the 's' in separators
		l, p := JumpForward(b, c)

		AssertLinePos(t, 2, 2, l, p) // expect space to be skipped
	})
	t.Run("Jump URL", func(t *testing.T) {
		c := b.NewCursor(3, 0)
		l, p := JumpForward(b, c)

		AssertLinePos(t, 3, 5, l, p) // expect space to be skipped
	})
	t.Run("Jump URL 2", func(t *testing.T) {
		c := b.NewCursor(3, 5) // where we ended the previous test
		l, p := JumpForward(b, c)

		AssertLinePos(t, 3, 8, l, p)
	})
	t.Run("Jump next line word continues", func(t *testing.T) {
		b := novi.BuildBuffer("foo bar", "bla 123")
		c := b.NewCursor(0, 4) // where we ended the previous test
		l, p := JumpForward(b, c)

		AssertLinePos(t, 1, 0, l, p)
	})
}

// Test 'b' behaviour
func TestJumpBackward(t *testing.T) {
	b := novi.BuildBuffer("This..isa line/;-with? separators", "", "  leading space",
		"https://github.com/some/repo.git?foo=a", "last line")

	t.Run("Find first from end", func(t *testing.T) {
		c := b.NewCursor(4, 8)
		l, p := JumpBackward(b, c)

		AssertLinePos(t, 4, 5, l, p)
	})
	t.Run("Find first from interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 16) // - before with
		l, p := JumpBackward(b, c)

		AssertLinePos(t, 0, 14, l, p)
	})
	t.Run("Find previous word", func(t *testing.T) {
		c := b.NewCursor(0, 14) // / after line
		l, p := JumpBackward(b, c)

		AssertLinePos(t, 0, 10, l, p)
	})
	t.Run("Match on previous line", func(t *testing.T) {
		c := b.NewCursor(4, 0)
		l, p := JumpBackward(b, c)

		AssertLinePos(t, 3, 37, l, p)
	})
}

func TestJumpForwardEnd(t *testing.T) {
	b := novi.BuildBuffer("This..isa line/;-with? separators", "", "  leading space",
		"https://github.com/some/repo.git?foo=a", "last line")

	t.Run("Find first from start", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		l, p := JumpForwardEnd(b, c)

		AssertLinePos(t, 0, 3, l, p)
	})
	t.Run("Find first from interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 3)
		l, p := JumpForwardEnd(b, c)

		AssertLinePos(t, 0, 5, l, p)
	})
	t.Run("Find next line URL", func(t *testing.T) {
		c := b.NewCursor(2, 14) // end line 3
		l, p := JumpForwardEnd(b, c)

		AssertLinePos(t, 3, 4, l, p) // end of https
	})
}

func TestJumpWordForwardEnd(t *testing.T) {
	b := novi.BuildBuffer("This..isa line/;-with? separators", "", "  leading space",
		"https://github.com/some/repo.git?foo=a", "last line")

	t.Run("Find first from start", func(t *testing.T) {
		c := b.NewCursor(0, 0)
		l, p := JumpWordForwardEnd(b, c)

		AssertLinePos(t, 0, 8, l, p)
	})
	t.Run("Find first from interpunction", func(t *testing.T) {
		c := b.NewCursor(0, 8)
		l, p := JumpWordForwardEnd(b, c)

		AssertLinePos(t, 0, 21, l, p) // ? after with
	})
	t.Run("Find next line URL", func(t *testing.T) {
		c := b.NewCursor(2, 14) // end line 3
		l, p := JumpWordForwardEnd(b, c)

		AssertLinePos(t, 3, 37, l, p) // end of URL
	})
}

func AssertWordMatches(t *testing.T, m []int, exp []int) {
	t.Helper()

	if len(m) != len(exp) {
		t.Fatalf("Didn't get equal sized expected/actual: %v - %v", m, exp)
	}

	for i, e := range m {
		if e != exp[i] {
			t.Errorf("Difference at position %d: got %d but expected %d", i, e, exp[i])
		}
	}
}
func TestWordStarts(t *testing.T) {
	// Words are sequences of alphanum *or* sequences of separators
	t.Run("Empty line", func(t *testing.T) {
		res := WordStarts(novi.Line{}, false)

		AssertWordMatches(t, res, []int{0})
	})
	t.Run("Expect nothing on all spaces", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("    ")), false)

		AssertWordMatches(t, res, []int{})
	})
	t.Run("Some simple words, variable spaces", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("  this   is  a     test")), false)

		AssertWordMatches(t, res, []int{2, 9, 13, 19})
	})
	t.Run("Mix of alphanum, separator words", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("this, is! a *!@&#^ test")), false)

		AssertWordMatches(t, res, []int{0, 4, 6, 8, 10, 12, 19})
	})
	t.Run("A URL", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("https://www.github.com/sample/repo.git?foo=bar")), false)
		AssertWordMatches(t, res, []int{0, 5, 8, 11, 12, 18, 19, 22, 23, 29, 30, 34, 35, 38, 39, 42, 43})
	})

	// Words are sequences of alphanum *or* separators
	t.Run("Empty line (same)", func(t *testing.T) {
		res := WordStarts(novi.Line{}, true)

		AssertWordMatches(t, res, []int{0})
	})
	t.Run("Expect nothing on all spaces (same)", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("    ")), true)

		AssertWordMatches(t, res, []int{})
	})
	t.Run("Some simple words, variable spaces (same)", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("  this   is  a     test")), true)

		AssertWordMatches(t, res, []int{2, 9, 13, 19})
	})
	t.Run("Mix of alphanum, separator words (same)", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("this, is! a *!@&#^ !test!")), true)

		AssertWordMatches(t, res, []int{0, 6, 10, 12, 19})
	})
	t.Run("A URL (same)", func(t *testing.T) {
		res := WordStarts(novi.Line([]rune("https://www.github.com/sample/repo.git?foo=bar")), true)
		AssertWordMatches(t, res, []int{0})
	})
}

func TestWordEnds(t *testing.T) {
	// Treat alnum/separators separately
	t.Run("Empty line", func(t *testing.T) {
		res := WordEnds(novi.Line{}, false)

		AssertWordMatches(t, res, []int{0})
	})
	t.Run("Expect nothing on all spaces", func(t *testing.T) {
		res := WordEnds(novi.Line([]rune("    ")), false)

		AssertWordMatches(t, res, []int{})
	})
	t.Run("Simple test on letters and spaces", func(t *testing.T) {
		res := WordEnds(novi.Line([]rune(" this   is  a     test ")), false)
		AssertWordMatches(t, res, []int{4, 9, 12, 21})
	})
	t.Run("Mix of alphanum, separator words", func(t *testing.T) {
		res := WordEnds(novi.Line([]rune("this, is! a *!@&#^ test")), false)

		AssertWordMatches(t, res, []int{3, 4, 7, 8, 10, 17, 22})
	})
	t.Run("A URL", func(t *testing.T) {
		res := WordEnds(novi.Line([]rune("https://www.github.com/sample/repo.git?foo=bar")), false)
		AssertWordMatches(t, res, []int{4, 7, 10, 11, 17, 18, 21, 22, 28, 29, 33, 34, 37, 38, 41, 42, 45})
	})

	// Treat alnum/separators as the same for defining words
	t.Run("Empty line (same)", func(t *testing.T) {
		res := WordEnds(novi.Line{}, true)

		AssertWordMatches(t, res, []int{0})
	})
	t.Run("Expect nothing on all spaces (same)", func(t *testing.T) {
		res := WordEnds(novi.Line([]rune("    ")), true)

		AssertWordMatches(t, res, []int{})
	})
	t.Run("Mix of alphanum, separator words (same)", func(t *testing.T) {
		res := WordEnds(novi.Line([]rune("this, is! a *!@&#^ !test!")), true)

		AssertWordMatches(t, res, []int{4, 8, 10, 17, 24})
	})
}
