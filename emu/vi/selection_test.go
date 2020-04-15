package viemu

import (
	"testing"

	"github.com/iivvoo/novi/novi"
)

func TestBlockSelection(t *testing.T) {
	t.Run("Test regular block delete", func(t *testing.T) {
		vi, _ := SetupViAndCursor(ModeSelect, 0, 4,
			"1234567890",
			"1234567890",
			"1234567890",
			"1234567890",
		)
		vi.HandleSelectionBlock(nil)
		vi.SelectionStart = *novi.NewCursor(vi.Editor.Buffer, 1, 3)
		vi.SelectionEnd = *novi.NewCursor(vi.Editor.Buffer, 3, 7)
		vi.HandleSelectRemove(nil)

		novi.AssertBufferMatch(t, vi.Editor.Buffer,
			"1234567890",
			"12390",
			"12390",
			"12390",
		)
	})
	t.Run("Test start/end swapped", func(t *testing.T) {
		vi, _ := SetupViAndCursor(ModeSelect, 0, 4,
			"1234567890",
			"1234567890",
			"1234567890",
			"1234567890",
		)
		vi.HandleSelectionBlock(nil)
		vi.SelectionStart = *novi.NewCursor(vi.Editor.Buffer, 3, 7)
		vi.SelectionEnd = *novi.NewCursor(vi.Editor.Buffer, 1, 3)
		vi.HandleSelectRemove(nil)

		novi.AssertBufferMatch(t, vi.Editor.Buffer,
			"1234567890",
			"12390",
			"12390",
			"12390",
		)
	})
	t.Run("Test different line lengths", func(t *testing.T) {
		vi, _ := SetupViAndCursor(ModeSelect, 0, 4,
			"1234567890",
			"12",
			"123",
			"1234567890",
		)
		vi.HandleSelectionBlock(nil)
		vi.SelectionStart = *novi.NewCursor(vi.Editor.Buffer, 1, 3)
		vi.SelectionEnd = *novi.NewCursor(vi.Editor.Buffer, 3, 7)
		vi.HandleSelectRemove(nil)

		novi.AssertBufferMatch(t, vi.Editor.Buffer,
			"1234567890",
			"12",
			"123",
			"12390",
		)
	})
}
