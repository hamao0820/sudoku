package main

import (
	"fmt"
	"image"
	"math"

	"github.com/hamao0820/sudoku/detect"
	"gocv.io/x/gocv"
)

func main() {
	src := gocv.IMRead("sample3.png", gocv.IMReadColor)
	defer src.Close()

	detect.FitSize(&src, 500, 500)

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

	edge := gocv.NewMat()
	defer edge.Close()
	gocv.Canny(gray, &edge, 50, 150)

	hierarchy := gocv.NewMat()
	defer hierarchy.Close()
	contours := gocv.FindContoursWithParams(edge, &hierarchy, gocv.RetrievalList, gocv.ChainApproxSimple)

	min_area := float64(src.Rows()*src.Cols()) * 0.2
	largeContors := gocv.NewPointsVector()
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		if gocv.ContourArea(contour) > min_area {
			largeContors.Append(contour)
		}
	}

	convexes := make([]gocv.Mat, largeContors.Size())
	for i := 0; i < largeContors.Size(); i++ {
		contour := largeContors.At(i)

		convex := gocv.NewMat()
		defer convex.Close()
		gocv.ConvexHull(contour, &convex, true, false) // convexはcontourの部分集合のインデックス
		convexes[i] = convex
	}

	polies := gocv.NewPointsVector()
	for i := 0; i < largeContors.Size(); i++ {
		contour := largeContors.At(i)
		converted := gocv.NewMat()
		defer converted.Close()
		convex := convexes[i]
		convexContour := gocv.NewPointVector()
		for j := 0; j < convex.Rows(); j++ {
			convexContour.Append(contour.At(int(convex.GetIntAt(j, 0))))
		}

		arcLen := gocv.ArcLength(convexContour, true)
		poly := gocv.ApproxPolyDP(convexContour, 0.02*arcLen, true)
		polies.Append(poly)
	}

	// 正方形に近いものを選ぶ
	selectedIndex := 0 // 選ばれたindex
	pMinError := math.MaxFloat64
	for i := 0; i < polies.Size(); i++ {
		poly := polies.At(i)
		if poly.Size() < 4 {
			continue
		}

		// polyの中から4つ選ぶ
		// 4つの点の組み合わせを全て試す
		indices := make([]int, poly.Size())
		for j := 0; j < poly.Size(); j++ {
			indices[j] = j
		}
		errors := make([]float64, 0) // 4つの点の組み合わせごとの誤差
		combs := make([][]int, 0)    // 4つの点の組み合わせ
		for comb := range combinations(indices, 4, 1) {
			combs = append(combs, comb)
		}
		for _, comb := range combs {
			quadrilateral := gocv.NewPointVector()
			for _, c := range comb {
				quadrilateral.Append(poly.At(c))
			}

			lineLens := make([]float64, 4)
			for j := 0; j < 4; j++ {
				lineLens[j] = math.Sqrt(math.Pow(float64(quadrilateral.At(j).X-quadrilateral.At((j+1)%4).X), 2) + math.Pow(float64(quadrilateral.At(j).Y-quadrilateral.At((j+1)%4).Y), 2))
			}
			base := math.Sqrt(gocv.ContourArea(quadrilateral))
			score := 0.0
			for _, l := range lineLens {
				score += math.Pow(math.Abs(1-l/base), 2)
			}
			errors = append(errors, score)
		}

		minError := errors[0] // 4つの点の組み合わせごとの誤差の最小値
		for _, e := range errors {
			if e < minError {
				minError = e
			}
		}

		if minError < pMinError {
			pMinError = minError
			selectedIndex = i
		}
	}

	poly := fixClockwise(polies.At(selectedIndex))

	// gocv.DrawContoursWithParams(&src, polies, selectedIndex, color.RGBA{255, 0, 0, 255}, 1, gocv.LineAA, hierarchy, 0, image.Pt(0, 0))

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

	quadrilateral := gocv.NewMat()
	defer quadrilateral.Close()

	warp := gocv.NewMat()
	defer warp.Close()

	fmt.Println(poly.ToPoints())
	gocv.WarpPerspective(src, &warp, gocv.GetPerspectiveTransform(poly, dst), image.Pt(src.Cols(), src.Rows()))

	trimmed := warp.Region(image.Rect(0, 0, size, size))
	defer trimmed.Close()

	gocv.Resize(trimmed, &trimmed, image.Pt(500, 500), 0, 0, gocv.InterpolationCubic)

	win := gocv.NewWindow("sudoku")
	defer win.Close()
	win.IMShow(trimmed)
	win.WaitKey(0)
}

func fixClockwise(poly gocv.PointVector) gocv.PointVector {
	points := poly.ToPoints()
	// 重心を求める
	center := image.Pt(0, 0)
	for _, p := range points {
		center.X += p.X
		center.Y += p.Y
	}
	center.X /= len(points)
	center.Y /= len(points)

	// 重心からの角度を求める
	angles := make([]float64, len(points))
	for i, p := range points {
		angles[i] = math.Atan2(float64(p.Y-center.Y), float64(p.X-center.X))
	}

	// 角度でソート
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			if angles[i] < angles[j] {
				angles[i], angles[j] = angles[j], angles[i]
				points[i], points[j] = points[j], points[i]
			}
		}
	}

	points = append(points[3:], points[:3]...) // 0, 1, 2, 3 -> 3, 0, 1, 2

	return gocv.NewPointVectorFromPoints(points)
}

func combinations(list []int, choose, buf int) (c chan []int) {
	c = make(chan []int, buf)
	go func() {
		defer close(c)
		switch {
		case choose == 0:
			c <- []int{}
		case choose == len(list):
			c <- list
		case len(list) < choose:
			return
		default:
			for i := 0; i < len(list); i++ {
				for subComb := range combinations(list[i+1:], choose-1, buf) {
					c <- append([]int{list[i]}, subComb...)
				}
			}
		}
	}()
	return
}
