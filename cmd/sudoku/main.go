package main

import (
	"fmt"
	"os"

	"github.com/gdey/sudoku"
)

func main() {
	var filename = "sample-puzzles/hard_1_puzzle.txt"

	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	b, err := sudoku.LoadFromFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	if err = b.Solve(); err != nil {
		panic(err)
	}
	fmt.Println(b)
}
