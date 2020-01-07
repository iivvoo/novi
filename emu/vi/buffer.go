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

// JumpForward jumps to the next sequence of alphanum or separators, skipping whitespace
func JumpForward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	line, pos := c.Line, c.Pos

	runeType := TypeUnknown
	for line < b.Length() {
		l := b.Lines[line]
		for pos < len(l) {
			cc := l[pos]

			newType := GetRuneType(cc)
			if runeType != TypeUnknown && newType != TypeSpace && newType != runeType {
				return line, pos
			}
			runeType = newType

			pos++
		}
		// treat the end-of-line as space since it separates words
		runeType = TypeSpace
		pos = 0
		line++
		// an empty line also matches
		if line < b.Length() && len(b.Lines[line]) == 0 {
			return line, pos
		}
	}
	return -1, -1
}

// JumpWordForward implements "W" behaviour, the begining of the next word
func JumpWordForward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	// cursor does not have to be bound to buffer

	line, pos := c.Line, c.Pos

	sepFound := false
	for line < b.Length() {
		l := b.Lines[line]
		for pos < len(l) {
			cc := l[pos]
			if unicode.IsSpace(cc) {
				sepFound = true
			} else if sepFound {
				return line, pos
			}
			pos++
		}
		sepFound = true
		// did we advance to a completely empty line? Then that's also a valid match
		line++
		pos = 0
		if line < b.Length() && len(b.Lines[line]) == 0 {
			return line, 0
		}
	}
	// return last character, even if it's space
	return b.Length() - 1, len(b.Lines[b.Length()-1]) - 1
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
