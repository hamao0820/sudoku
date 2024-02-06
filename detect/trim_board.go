package detect

import (
	"fmt"
	"image"
	"math"

	"gocv.io/x/gocv"
)

func GetCells(img gocv.Mat) [][]gocv.Mat {
	dx := float64(img.Cols()) / 9
	dy := float64(img.Rows()) / 9
	cells := make([][]gocv.Mat, 9)
	for y := 0; y < 9; y++ {
		cells[y] = make([]gocv.Mat, 9)
		for x := 0; x < 9; x++ {
			cells[y][x] = img.Region(image.Rect(int(float64(x)*dx), int(float64(y)*dy), int(float64(x+1)*dx), int(float64(y+1)*dy)))
		}
	}
	return cells
}

func GetSquare(img gocv.Mat) (gocv.Mat, error) {
	FitSize(&img, 500, 500)

	gray := ToGray(img)
	defer gray.Close()

	edge := FindEdge(gray)
	defer edge.Close()

	contours, _ := FindContours(edge)
	defer contours.Close()

	min_area := float64(img.Rows()*img.Cols()) * 0.2
	largeContours := FilterContours(contours, min_area)

	convexes := GetConvexes(largeContours)

	polies := GetPolygons(largeContours, convexes)
	defer polies.Close()

	// 正方形に近いものを選ぶ
	selectedIndex, _ := SelectNearestSquareIndex(polies)

	if selectedIndex == -1 {
		return gocv.NewMat(), fmt.Errorf("not found")
	}

	poly := fixClockwise(polies.At(selectedIndex))

	warp, size := getSquareWarpPerspectiveTransformed(img, poly)

	square := warp.Region(image.Rect(0, 0, size, size))
	return square, nil
}

func getSquarePerspectiveTransform(poly gocv.PointVector) (gocv.Mat, int) {
	dst := gocv.NewPointVector()
	defer dst.Close()

	lineLensSum := 0.0
	for j := 0; j < 4; j++ {
		lineLensSum += math.Sqrt(math.Pow(float64(poly.At(j).X-poly.At((j+1)%4).X), 2) + math.Pow(float64(poly.At(j).Y-poly.At((j+1)%4).Y), 2))
	}
	size := int(lineLensSum / 4)
	dst.Append(image.Pt(0, 0))
	dst.Append(image.Pt(0, size))
	dst.Append(image.Pt(size, size))
	dst.Append(image.Pt(size, 0))

	return gocv.GetPerspectiveTransform(poly, dst), size
}

func getSquareWarpPerspectiveTransformed(img gocv.Mat, poly gocv.PointVector) (gocv.Mat, int) {
	perspective, size := getSquarePerspectiveTransform(poly)
	defer perspective.Close()
	warp := gocv.NewMat()

	gocv.WarpPerspective(img, &warp, perspective, image.Pt(img.Cols(), img.Rows()))
	return warp, size
}
