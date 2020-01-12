package viemu

import (
	"fmt"
	"unicode"

	"github.com/iivvoo/ovim/ovim"
)

/*
 * Contains vi-specific buffer operations, manipulations. If reusable, move to ovim
 */
/*
  w 	Move to next word
  W 	Move to next blank delimited word
  b 	Move to the beginning of the word
  B 	Move to the beginning of blank delimited word
*/

// RuneType identifies the type of a rune
type RuneType int

// The possible values a RuneType can have
const (
	TypeAlNum RuneType = iota
	TypeSep
	TypeSpace
	TypeUnknown
)

// GetRuneType returns if the rune is alphanum, whitespace or separator
func GetRuneType(r rune) RuneType {
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		return TypeAlNum
	}
	if unicode.IsSpace(r) {
		return TypeSpace
	}
	return TypeSep
}

// WordStarts Finds al the "word starts" in a given line
func WordStarts(l ovim.Line, SepAndAlnum bool) []int {
	if len(l) == 0 {
		return []int{0}
	}

	var res []int

	prevRune := TypeUnknown
	for pos, c := range l {
		newRune := GetRuneType(c)

		if SepAndAlnum && newRune == TypeSep {
			newRune = TypeAlNum
		}
		if prevRune != newRune && newRune != TypeSpace {
			res = append(res, pos)
		}
		prevRune = newRune
	}

	return res
}

// WordEnds finds all the endings of words in a given line.
// words are either sequences of alphanum or non-whitespace,
// e,g ab123 []=-- or sequences of alphanum *and* non-ws
// This. is a?!@ sentence...
//    ^^  ^
func WordEnds(l ovim.Line, SepAndAlnum bool) []int {
	if len(l) == 0 {
		return []int{0}
	}

	var res []int

	prevRune := TypeUnknown
	for i := len(l) - 1; i >= 0; i-- {
		newRune := GetRuneType(l[i])

		// treat separators like alnum?
		if SepAndAlnum && newRune == TypeSep {
			newRune = TypeAlNum
		}
		if prevRune != newRune && newRune != TypeSpace {
			res = append(res, i)
		}
		prevRune = newRune
	}
	return res
}

// JumpForward jumps to the next sequence of alphanum or separators, skipping whitespace
func JumpForward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	line, pos := c.Line, c.Pos

	for line < b.Length() {
		l := b.Lines[line]

		positions := WordStarts(l, false)

		for _, p := range positions {
			if p > pos {
				return line, p
			}
		}
		pos = -1 // make sure we're really smaller, so we will match on pos 0
		line++
	}
	return b.Length() - 1, len(b.Lines[b.Length()-1]) - 1
}

// JumpWordForward implements "W" behaviour, the begining of the next word
func JumpWordForward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	line, pos := c.Line, c.Pos

	for line < b.Length() {
		l := b.Lines[line]

		positions := WordStarts(l, true)

		for _, p := range positions {
			if p > pos {
				return line, p
			}
		}
		pos = -1 // make sure we're really smaller, so we will match on pos 0
		line++
	}
	// jump to the very last character in the buffer
	return b.Length() - 1, len(b.Lines[b.Length()-1]) - 1
}

// JumpBackward implements "b" behaviour, the beginning of the previous sequence of alphanum or other non-whitespace
func JumpBackward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	line, pos := c.Line, c.Pos

	for line >= 0 {
		l := b.Lines[line]
		positions := WordStarts(l, false)
		lastPos := -1

		for _, p := range positions {
			// There must be a previous pos to return, the current one must be larger,
			// and the previous should be smaller (it could be equal!)
			if lastPos != -1 && p >= pos && lastPos < pos {
				return line, lastPos
			}
			lastPos = p
		}
		// if all matches were smaller than pos, return the last
		if lastPos != -1 && lastPos < pos {
			return line, lastPos
		}

		// continue to the next line, position cursor at the end
		line--
		if line >= 0 {
			pos = len(b.Lines[line]) + 1 // add one so we're larger than a match at the end
		}
	}
	return 0, 0
}

// JumpWordBackward implements "B" behaviour, the beginning of the previous word, skipping everything non-alphanum
func JumpWordBackward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	line, pos := c.Line, c.Pos
	lastLine, lastPos := line, pos

	wordFound := false
	didMove := false

	for line >= 0 {
		l := b.Lines[line]
		// if we advanced at least one character and ended up on an empty line, we're good
		if didMove && len(l) == 0 {
			return line, 0
		}
		// pos is the position where we could insert, but it doesn't mean there's already a character there
		// mostly relevant on empty lines
		for pos >= 0 && pos < len(l) {
			cc := l[pos]
			if didMove && (unicode.IsLetter(cc) || unicode.IsNumber(cc)) {
				wordFound = true
			} else if wordFound {
				// we found a word, now a non-alphanum, so our desired position is the
				// previous
				return line, lastPos
			}
			lastPos = pos
			pos--
			didMove = true
		}
		// if we scanned letters but ended up on pos 0, that's also the start of a word
		if wordFound {
			// does this cover all cases?
			return line, 0
		}
		lastLine = line
		line--
		if line >= 0 {
			pos = len(b.Lines[line]) - 1
		}
		didMove = true
	}
	fmt.Println(lastLine)
	return 0, 0
}
