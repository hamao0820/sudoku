package main

import (
	"fmt"
	"time"

	"github.com/hamao0820/sudoku/detect"
	"github.com/hamao0820/sudoku/ocr"
	"gocv.io/x/gocv"
)

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("Test")
	defer window.Close()
	img := gocv.NewMat()
	defer img.Close()

	frame := 0
	cellDigits := make([][]int, 9)
	for {
		time.Sleep(100 * time.Millisecond)
		webcam.Read(&img)

		square, err := detect.GetSquare(img)
		if err == nil {
			frame++
			if frame >= 3 {
				cells := ocr.SplitCells(square)
				for y := 0; y < 9; y++ {
					cellDigits[y] = make([]int, 9)
					for x := 0; x < 9; x++ {
						cell := gocv.NewMat()
						gocv.CvtColor(cells[y][x], &cell, gocv.ColorBGRToGray)
						if cell.Empty() {
							continue
						}
						digit := ocr.OCR(cell)
						cellDigits[y][x] = digit
					}
				}
				break
			}
		} else {
			frame = 0
		}

		window.IMShow(img)
		window.WaitKey(1)
	}

	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			fmt.Printf("%d ", cellDigits[y][x])
		}
		fmt.Println()
	}
}
