package novi

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

func TestCursorValidate(t *testing.T) {
	t.Run("Just valid", func(t *testing.T) {
		c := BuildBuffer("line 1", "line 2").NewCursor(1, 3)
		if !c.Validate() {
			t.Errorf("Cursor did not validate but I did expect it to.")
		}
	})
	t.Run("Pos too far", func(t *testing.T) {
		c := BuildBuffer("line 1", "line 2").NewCursor(1, 30)
		if c.Validate() {
			t.Errorf("Cursor did validate but I did expect it not to.")
		}
		if c.Line != 1 || c.Pos != 5 {
			t.Errorf("Expected cursor to be placed at (1, 5) but got (%d, %d)", c.Line, c.Pos)
		}
	})
	t.Run("Line too far", func(t *testing.T) {
		c := BuildBuffer("line 1", "line 2").NewCursor(10, 3)
		if c.Validate() {
			t.Errorf("Cursor did validate but I did expect it not to.")
		}
		if c.Line != 1 || c.Pos != 3 {
			t.Errorf("Expected cursor to be placed at (1, 3) but got (%d, %d)", c.Line, c.Pos)
		}
	})
	t.Run("All too far", func(t *testing.T) {
		c := BuildBuffer("line 1", "line 2").NewCursor(10, 30)
		if c.Validate() {
			t.Errorf("Cursor did validate but I did expect it not to.")
		}
		if c.Line != 1 || c.Pos != 5 {
			t.Errorf("Expected cursor to be placed at (1, 5) but got (%d, %d)", c.Line, c.Pos)
		}
	})
	t.Run("Empty Buffer", func(t *testing.T) {
		c := BuildBuffer().NewCursor(10, 30)
		if c.Validate() {
			t.Errorf("Cursor did validate but I did expect it not to.")
		}
		if c.Line != 0 || c.Pos != 0 {
			t.Errorf("Expected cursor to be placed at (0, 0) but got (%d, %d)", c.Line, c.Pos)
		}
	})
}
