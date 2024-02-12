package main

import (
	"fmt"
	"time"

	"github.com/hamao0820/sudoku/detect"
	"github.com/hamao0820/sudoku/ocr"
	"github.com/hamao0820/sudoku/solver"
	"gocv.io/x/gocv"
)

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("Test")
	defer window.Close()
	original := gocv.NewMat()
	defer original.Close()

	frame := 0
	cellDigits := [9][9]int{}
	for {
		time.Sleep(100 * time.Millisecond)
		webcam.Read(&original)

		display := original.Clone()
		if detect.DrawSquare(&display) {
			square, err := detect.GetSquare(original)
			if err == nil {
				frame++
				if frame >= 3 {
					cells := detect.SplitCell(square)
					for y := 0; y < 9; y++ {
						cellDigits[y] = [9]int{}
						for x := 0; x < 9; x++ {
							digit := ocr.OCR(cells[y][x])
							cellDigits[y][x] = digit
						}
					}
					break
				}
			} else {
				frame = 0
			}
		}

		window.IMShow(display)
		window.WaitKey(1)
	}

	b := solver.Board(cellDigits)
	if !b.Verify() {
		fmt.Println("Invalid board")
	}

	fmt.Println("Detected digits:")
	fmt.Println(b.String())

	solver.Backtrack(&b)

	fmt.Println("Solved board:")
	fmt.Println(b.String())

}
