package viemu

import (
	"strings"

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
/*
 * de 'w' variant sprint woord of interpunctie
 * dit...iseen,?rare zin!?_bla:la
 *  | |  |    | |    |  | |   ||
 */

type RuneType int

const (
	TypeWord RuneType = iota
	TypeSep
	TypeSpace
	TypeUnknown
)

func GetRuneType(r rune) RuneType {
	// golang probable has better (unicode) tools for this
	if strings.IndexRune(".,;'/", r) != -1 {
		return TypeSep
	}
	if strings.IndexRune(" \n\t", r) != -1 {
		return TypeSpace
	}
	return TypeWord
}

func JumpForward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	// basically make a distinction between separator-words, number/letter words, whitespace
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
		pos = 0
		line++
		// an empty line also matches
		if line < b.Length() && len(b.Lines[line]) == 0 {
			return line, pos
		}
	}
	return -1, -1
}

// JumpWordForward implements "W" behaviour, the begining of the next
// word
func JumpWordForward(b *ovim.Buffer, c *ovim.Cursor) (int, int) {
	// cursor does not have to be bound to buffer

	line, pos := c.Line, c.Pos

	sepFound := false
	for line < b.Length() {
		l := b.Lines[line]
		for pos < len(l) {
			cc := l[pos]
			if cc == ' ' { // XXX tab?
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
			return line, pos
		}
	}
	// return last character, even if it's space
	return b.Length() - 1, len(b.Lines[b.Length()-1]) - 1
}
