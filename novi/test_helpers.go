package novi

import "testing"

/*
 * Test helpers that we want to share with external packages
 */

// AssertCursor asserts the cursor matches the expected line, pos
func AssertCursor(t *testing.T, c *Cursor, line, pos int) {
	t.Helper()

	if c.Line != line || c.Pos != pos {
		t.Errorf("Cursor (%d, %d) didn't match what was expected (%d, %d)",
			c.Line, c.Pos, line, pos)
	}
}

// AssertBufferMatch asserts the buffer matches the expected string slice
func AssertBufferMatch(t *testing.T, b *Buffer, expected ...string) {
	t.Helper()

	if a, b := b.Length(), len(expected); a != b {
		t.Errorf("Buffer size %d does not match expected size %d", a, b)
	}
	for i, e := range expected {
		if i >= b.Length() {
			break
		}
		if got := string(b.Lines[i]); e != got {
			t.Errorf("First mismatch at line %d\nexpected: %s\ngot     : %s", i, e, got)
		}
	}
}

// AssertBufferModified asserts the buffer was modified
func AssertBufferModified(t *testing.T, b *Buffer, modified bool) {
	t.Helper()

	if modified && !b.Modified {
		t.Errorf("Expected buffer to be modified, but it wasn't")
	} else if !modified && b.Modified {
		t.Errorf("Expected buffer NOT to be modified, but it actually WAS")
	}
}

// BuildBuffer creates a new buffer based on the supplied strings
func BuildBuffer(lines ...string) *Buffer {
	return NewBuffer().LoadStrings(lines)
}
