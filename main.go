package main

import (
	"fmt"

	"github.com/hamao0820/sudoku/detect"
	"gocv.io/x/gocv"
)

func main() {
	src := gocv.IMRead("sample2.png", gocv.IMReadColor)
	defer src.Close()

	square, err := detect.GetSquare(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	win := gocv.NewWindow("sudoku")
	defer win.Close()
	win.IMShow(square)
	win.WaitKey(0)
}
