package main

import (
	"fmt"

	"github.com/hamao0820/sudoku/solver"
)

func main() {
	b := solver.NewBoard()
	err := b.FromFile("solver/test/easy1.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(b.String())
}
