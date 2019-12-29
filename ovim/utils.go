package ovim

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
)

// RecoverFromPanic attempts to print a usable stacktrace in a panic-recover
func RecoverFromPanic(cleanup func()) {
	/*
			A normal debug.Stack() at this point looks like this:
			goroutine 23 [running]:
				runtime/debug.Stack(0xc00012a000, 0x584f40, 0xc0000d0120)
		       /opt/go1.13.1/src/runtime/debug/stack.go:24 +0x9d
				main.start.func1.1(0xc0000882a0, 0x5c3de0, 0xc00012a000)
		       /projects/ovim/cmd/ovim/main.go:65 +0x8d
				panic(0x584f40, 0xc0000d0120)
		       /opt/go1.13.1/src/runtime/panic.go:679 +0x1b2
				gitlab.com/iivvoo/ovim/ovim.(*TermUI).RenderTerm(0xc0000b8700)
		       /projects/ovim/ovim/term.go:90 +0x487
				main.start.func1(0xc0000882a0, 0x5c3de0, 0xc00012a000, 0xc00008c720, 0xc0000b8700)
		       /projects/ovim/cmd/ovim/main.go:110 +0x229
				created by main.start

			We don't care about the trace up until the panic and the line after it,
			so we need to do some fancy counting to skip that. Also, it's nice to
			have the actual error at the bottom
	*/
	if r := recover(); r != nil {
		cleanup()
		s := strings.Split(string(debug.Stack()), "\n")
		panicCheck := 0
		fmt.Println(s[0])
		for _, entry := range s {
			if panicCheck == 2 {
				fmt.Println(entry)
			} else if strings.HasPrefix(entry, "panic(") {
				panicCheck = 1
			} else if panicCheck == 1 {
				panicCheck = 2
			}
		}
		fmt.Printf("%s\n", r)
	}
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer sourceFile.Close()

	newFile, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	return err
}
