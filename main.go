package main

import "os"

func main() {
	(&Harry{
		MakeArgs: os.Args[1:],
	}).MakeMyDay()
}
