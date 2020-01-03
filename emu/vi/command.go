package viemu

import (
	"math"
	"strconv"
	"strings"
)

// ParseCommand attempts to parse vi command mode commands into count + command
func ParseCommand(command string) (int, string) {
	/*
	 * a vi(m?) command has the structure
	 * <number?>character
	 * <number?>character(<number?>character)? e.g. 2d3d -> 6dd, or d10d -> 10dd
	 *
	 * (vim actually understands <num><keyup>!)
	 * Also, 2d0 is executed immediately and deletes to home
	 */
	count := 1
	cmd := ""
	numMode := false

	storeCount := func(sub string) {
		if sub != "" {
			if v, e := strconv.Atoi(sub); e != nil {
				// can only happen with an overflowing number. Not sure what the right
				// cause of action would be
				count = math.MaxInt64
			} else {
				count *= v
			}
		}
	}

	storeCommand := func(sub string) {
		if sub != "" {
			cmd += sub
		}
	}

	s, e := 0, 0
	for e < len(command) {
		r := rune(command[e])
		if strings.IndexRune("0123456789", r) != -1 {
			if numMode {
				e++
				continue
			}
			numMode = true
			storeCommand(command[s:e])
			s = e
		} else {
			if !numMode {
				e++
				continue
			}
			numMode = false
			storeCount(command[s:e])
			s = e
		}
		e++
	}
	// store whatever is left
	if numMode {
		storeCount(command[s:e])
	} else {
		storeCommand(command[s:e])
	}

	return count, cmd
}
