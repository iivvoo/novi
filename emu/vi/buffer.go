package viemu

import (
	"unicode"

	"github.com/iivvoo/novi/novi"
)

/*
 * Contains vi-specific buffer operations, manipulations. If reusable, move to novi
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
func WordStarts(l novi.Line, SepAndAlnum bool) []int {
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
func WordEnds(l novi.Line, SepAndAlnum bool) []int {
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
			// reverse append
			res = append([]int{i}, res...)
		}
		prevRune = newRune
	}
	return res
}

// JumpAlNumSepForward jumps to the start of the next word. alnumSepSame defines if
// words include separators, or if they count as separate words
func JumpAlNumSepForward(b *novi.Buffer, c *novi.Cursor, alnumSepSame bool) (int, int) {
	line, pos := c.Line, c.Pos

	for line < b.Length() {
		l := b.Lines[line]

		positions := WordStarts(l, alnumSepSame)

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

// JumpForward jumps to the next sequence of alphanum or separators, skipping whitespace
func JumpForward(b *novi.Buffer, c *novi.Cursor) (int, int) {
	return JumpAlNumSepForward(b, c, false)
}

// JumpWordForward implements "W" behaviour, the begining of the next word
func JumpWordForward(b *novi.Buffer, c *novi.Cursor) (int, int) {
	return JumpAlNumSepForward(b, c, true)
}

// JumpAlNumSepBackward jumps to the start of the previous word. alnumSepSame defines if
// words include separators, or if they count as separate words
func JumpAlNumSepBackward(b *novi.Buffer, c *novi.Cursor, alnumSepSame bool) (int, int) {
	line, pos := c.Line, c.Pos

	for line >= 0 {
		l := b.Lines[line]
		positions := WordStarts(l, alnumSepSame)
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

// JumpBackward implements "b" behaviour, the beginning of the previous sequence of alphanum or other non-whitespace
func JumpBackward(b *novi.Buffer, c *novi.Cursor) (int, int) {
	return JumpAlNumSepBackward(b, c, false)
}

// JumpWordBackward implements "B" behaviour, the beginning of the previous word, skipping everything non-alphanum
func JumpWordBackward(b *novi.Buffer, c *novi.Cursor) (int, int) {
	return JumpAlNumSepBackward(b, c, true)
}

// JumpAlNumSepForwardEnd jumps to the end of the next word. alnumSepSame defines if
// words include separators, or if they count as separate words
func JumpAlNumSepForwardEnd(b *novi.Buffer, c *novi.Cursor, alnumSepSame bool) (int, int) {
	line, pos := c.Line, c.Pos

	for line < b.Length() {
		l := b.Lines[line]

		positions := WordEnds(l, alnumSepSame)

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

// JumpForwardEnd implements "b" behaviour, the beginning of the previous sequence of alphanum or other non-whitespace
func JumpForwardEnd(b *novi.Buffer, c *novi.Cursor) (int, int) {
	return JumpAlNumSepForwardEnd(b, c, false)
}

// JumpWordForwardEnd implements "B" behaviour, the beginning of the previous word, skipping everything non-alphanum
func JumpWordForwardEnd(b *novi.Buffer, c *novi.Cursor) (int, int) {
	return JumpAlNumSepForwardEnd(b, c, true)
}
