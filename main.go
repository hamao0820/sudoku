package main

import (
	"fmt"

	"github.com/hamao0820/sudoku/solver"
)

func main() {
	b, err := solver.FromFile("solver/test/easy1.txt")
	fmt.Println(b.String())
	if err != nil {
		panic(err)
	}
	solver.Backtrack(b)
	fmt.Println(b.String())
}
