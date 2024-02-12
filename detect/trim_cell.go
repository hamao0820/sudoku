package detect

import (
	"image"

	"gocv.io/x/gocv"
)

func SplitCell(square gocv.Mat) [][]gocv.Mat {
	cells := make([][]gocv.Mat, 9)
	dx := float64(square.Cols()) / 9
	dy := float64(square.Rows()) / 9

	padding := 1.0
	for y := 0; y < 9; y++ {
		cells[y] = make([]gocv.Mat, 9)
		for x := 0; x < 9; x++ {
			cells[y][x] = square.Region(image.Rect(int(float64(x)*dx+padding), int(float64(y)*dy+padding), int(float64(x+1)*dx-padding), int(float64(y+1)*dy-padding)))
		}
	}
	return cells
}
