package ovim

import (
	"testing"
)

func TestCursorFiltering(t *testing.T) {
	t.Run("Filter Before, regular", func(t *testing.T) {
		cs := Cursors{&Cursor{Line: 10, Pos: 5}, &Cursor{Line: 9, Pos: 12}, &Cursor{Line: 12, Pos: 0}}
		before := cs.Before(cs[0])

		if len(before) != 1 {
			t.Fatalf("Expected exactly 1 cursor before given cursor, but found %d", len(before))
		}

		if before[0] != cs[1] {
			t.Fatalf("Expected cursor at (9,12) before but got %v", before[0])
		}
	})
	t.Run("Filter After, regular", func(t *testing.T) {
		cs := Cursors{&Cursor{Line: 10, Pos: 5}, &Cursor{Line: 9, Pos: 12}, &Cursor{Line: 12, Pos: 0}}
		after := cs.After(cs[0])

		if len(after) != 1 {
			t.Fatalf("Expected exactly 1 cursor after given cursor, but found %d", len(after))
		}

		if after[0] != cs[2] {
			t.Fatalf("Expected cursor at (12,0) after but got %v", after[0])
		}
	})
	t.Run("Filter Before, same line", func(t *testing.T) {
		cs := Cursors{&Cursor{Line: 10, Pos: 5}, &Cursor{Line: 10, Pos: 12}, &Cursor{Line: 10, Pos: 0}}
		before := cs.Before(cs[0])

		if len(before) != 1 {
			t.Fatalf("Expected exactly 1 cursor before given cursor, but found %d", len(before))
		}

		if before[0] != cs[2] {
			t.Fatalf("Expected cursor at (10,0) before but got %v", before[0])
		}
	})
	t.Run("Filter After, regular", func(t *testing.T) {
		cs := Cursors{&Cursor{Line: 10, Pos: 5}, &Cursor{Line: 10, Pos: 12}, &Cursor{Line: 10, Pos: 0}}
		after := cs.After(cs[0])

		if len(after) != 1 {
			t.Fatalf("Expected exactly 1 cursor after given cursor, but found %d", len(after))
		}

		if after[0] != cs[1] {
			t.Fatalf("Expected cursor at (10,12) after but got %v", after[0])
		}
	})
}
