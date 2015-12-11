package main

import "fmt"

func panicIfError(err error, message string) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}
