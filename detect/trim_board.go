package detect

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func DrawSquare(img *gocv.Mat) bool {
	FitSize(img, 500, 500)

	gray := ToGray(*img)
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
		return false
	}
	if polies.At(selectedIndex).Size() != 4 {
		return false
	}

	gocv.DrawContours(img, polies, selectedIndex, color.RGBA{255, 0, 0, 255}, 3)
	return true
}

func GetSquare(img gocv.Mat) (gocv.Mat, error) {
	FitSize(&img, 500, 500)

	filtered := BilateralFilter(img)
	defer filtered.Close()

	gray := ToGray(filtered)
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
	selectedIndex, errorScore := SelectNearestSquareIndex(polies)

	if selectedIndex == -1 {
		return gocv.NewMat(), fmt.Errorf("not found")
	}

	if errorScore > 0.1 {
		return gocv.NewMat(), fmt.Errorf("not square")
	}

	poly := fixClockwise(polies.At(selectedIndex))
	if poly.Size() != 4 {
		return gocv.NewMat(), fmt.Errorf("not square")
	}

	padding := 1
	points := poly.ToPoints()
	points[0].X = min(img.Cols(), points[0].X+padding)
	points[0].Y = min(img.Rows(), points[0].Y+padding)
	points[1].X = min(img.Cols(), points[1].X+padding)
	points[1].Y = max(0, points[1].Y-padding)
	points[2].X = max(0, points[2].X-padding)
	points[2].Y = max(0, points[2].Y-padding)
	points[3].X = max(0, points[3].X-padding)
	points[3].Y = min(img.Rows(), points[3].Y+padding)

	poly = gocv.NewPointVectorFromPoints(points)

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
