package main

import (
	"fmt"
	"os"
)

// CheckFatal : Check if err is set and print the error string
// then call CleanQuit to exit the program.
func CheckFatal(e error) {
	if e != nil {
		fmt.Fprint(os.Stderr, "Error: ", e)
		CleanQuit()
	}
}

// CleanQuit : Call SaveJSon and exit the program.
func CleanQuit() {
	//SaveJSon()
	os.Exit(0)
}
