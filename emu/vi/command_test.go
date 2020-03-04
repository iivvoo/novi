package viemu

import (
	"math"
	"testing"
)

func AssertCommand(t *testing.T, gotCount int, gotCmd string, expectCount int, expectCmd string) {
	t.Helper()

	if gotCount != expectCount {
		t.Errorf("Didn't get expected count %d, got %d", expectCount, gotCount)
	}
	if gotCmd != expectCmd {
		t.Errorf("Didn't get expected cmd %s, got %s", expectCmd, gotCmd)
	}
}

func TestCommandParsing(t *testing.T) {
	t.Run("Test single character", func(t *testing.T) {
		count, cmd := ParseCommand("l")

		AssertCommand(t, count, cmd, 1, "l")
	})
	t.Run("Test single with count", func(t *testing.T) {
		count, cmd := ParseCommand("32l")

		AssertCommand(t, count, cmd, 32, "l")

	})
	t.Run("Test double character", func(t *testing.T) {
		count, cmd := ParseCommand("dd")

		AssertCommand(t, count, cmd, 1, "dd")
	})
	t.Run("Test double character, weird mix", func(t *testing.T) {
		count, cmd := ParseCommand("2d3d")

		AssertCommand(t, count, cmd, 6, "dd")
	})
	t.Run("Test double character, count in between", func(t *testing.T) {
		count, cmd := ParseCommand("d3d")

		AssertCommand(t, count, cmd, 3, "dd")
	})
	t.Run("Test very long count", func(t *testing.T) {
		count, cmd := ParseCommand("99999999999999999999999999999999999999999999999999999999999999999999999l")

		AssertCommand(t, count, cmd, math.MaxInt64, "l")
	})
	t.Run("Test very long count 2", func(t *testing.T) {
		ParseCommand("99999999999999999999999999999999999999999999999999999999999999999999999d999d")

		// at this point I don't care about the count, as long as it doesn't crash
	})
	t.Run("Test just a number", func(t *testing.T) {
		count, cmd := ParseCommand("999")

		AssertCommand(t, count, cmd, 999, "")
	})
	t.Run("Test empty", func(t *testing.T) {
		count, cmd := ParseCommand("")

		AssertCommand(t, count, cmd, 1, "")
	})
}
