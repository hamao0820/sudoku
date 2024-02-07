package detect

import (
	"image"

	"gocv.io/x/gocv"
)

func GetCells(square gocv.Mat) [][]gocv.Mat {
	blurred := gocv.NewMat()
	defer blurred.Close()
	gocv.BilateralFilter(square, &blurred, 15, 150, 150)
	binary := gocv.NewMat()
	defer binary.Close()
	gocv.CvtColor(blurred, &binary, gocv.ColorBGRToGray)
	gocv.AdaptiveThreshold(binary, &binary, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 21, 0)

	dx := float64(binary.Cols()) / 9
	dy := float64(binary.Rows()) / 9
	cells := make([][]gocv.Mat, 9)
	padding := 3.0
	for y := 0; y < 9; y++ {
		cells[y] = make([]gocv.Mat, 9)
		for x := 0; x < 9; x++ {
			cells[y][x] = binary.Region(image.Rect(int(float64(x)*dx+padding), int(float64(y)*dy+padding), int(float64(x+1)*dx-padding), int(float64(y+1)*dy-padding)))
		}
	}

	return cells
}
