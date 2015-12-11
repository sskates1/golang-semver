package main

import (
	"fmt"
)

//SemVer struct
type SemVer struct {
	major      int
	minor      int
	patch      int
	preRelease string
	metadata   string
}

func main() {
	_ = getRegexes()
	fmt.Println("regexes loaded properly")
}
