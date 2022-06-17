//go:build windows
// +build windows

package colors

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

var Colors = map[string]string{
	"black":       "\u001b[30m",
	"red":         "\u001b[31m",
	"green":       "\u001b[32m",
	"yellow":      "\u001b[33m",
	"blue":        "\u001b[34m",
	"magenta":     "\u001b[35m",
	"cyan":        "\u001b[36m",
	"white":       "\u001b[37m",
	"reset":       "\u001b[0m",
	"brightblack": "\u001b[30;1m",
}

// Handle Virtual Termminal mode if on Windows
func init() {
	if runtime.GOOS == "windows" {
		h := syscall.Handle(os.Stdout.Fd())
		k32dll := syscall.NewLazyDLL("kernel32.dll")
		mode := k32dll.NewProc("SetConsoleMode")

		if _, _, err := mode.Call(uintptr(h), 0x001|0x002|0x004); err != nil && err.Error() != "The operation completed successfully." {
			for k := range Colors {
				Colors[k] = ""
			}
		}
	}
}

// Color functions.
func Red(msg any) string    { return fmt.Sprintf("%s%v%s", Colors["red"], msg, Colors["reset"]) }
func Green(msg any) string  { return fmt.Sprintf("%s%v%s", Colors["green"], msg, Colors["reset"]) }
func Yellow(msg any) string { return fmt.Sprintf("%s%v%s", Colors["yellow"], msg, Colors["reset"]) }
func Blue(msg any) string   { return fmt.Sprintf("%s%v%s", Colors["blue"], msg, Colors["reset"]) }
func Cyan(msg any) string   { return fmt.Sprintf("%s%v%s", Colors["cyan"], msg, Colors["reset"]) }
func Whtie(msg any) string  { return fmt.Sprintf("%s%v%s", Colors["white"], msg, Colors["reset"]) }
func Black(msg any) string  { return fmt.Sprintf("%s%v%s", Colors["black"], msg, Colors["reset"]) }
func Gray(msg any) string   { return fmt.Sprintf("%s%v%s", Colors["brightblack"], msg, Colors["reset"]) }
