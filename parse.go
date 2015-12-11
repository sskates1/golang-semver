package main

import (
	"fmt"
)

func parse(version string, loose string) *SemVer {
	fmt.Println(version, loose)
	x := new(SemVer)
	return x
}
