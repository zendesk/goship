package color

import (
	"fmt"
)

const colorGreen = "\x1b[32;1m"
const colorRed = "\x1b[31;1m"
const colorYellow = "\x1b[33;1m"
const colorReset = "\x1b[0m"

// PrintGreen prints string in green color
func PrintGreen(s string) {
	fmt.Printf("%s%s%s", colorGreen, s, colorReset)
}

// PrintRed prints string in red color
func PrintRed(s string) {
	fmt.Printf("%s%s%s", colorRed, s, colorReset)
}

// PrintYellow print string in yellow color
func PrintYellow(s string) {
	fmt.Printf("%s%s%s", colorYellow, s, colorReset)
}
