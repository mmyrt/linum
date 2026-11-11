package main

import (
	"fmt"
	"log"
	"os"
)

const (
	green = "\033[32m"
	blue  = "\033[34m"
	red   = "\033[31m"
	reset = "\033[0m"
)

func success(format string, args ...any) {
	fmt.Fprintf(os.Stdout, green+format+reset+"\n", args...)
}

func info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, blue+format+reset+"\n", args...)
}

func failure(format string, args ...any) {
	fmt.Fprintf(os.Stderr, red+format+reset+"\n", args...)
}

func fatal(args ...any) {
	log.Fatal(red + fmt.Sprint(args...) + reset)
}
