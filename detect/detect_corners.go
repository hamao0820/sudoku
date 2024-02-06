package detect

import (
	"image"
	"math"

	"gocv.io/x/gocv"
)

func FitSize(img *gocv.Mat, h, w int) {
	size := img.Size()
	f := math.Min(float64(h)/float64(size[0]), float64(w)/float64(size[1]))
	gocv.Resize(*img, img, image.Point{}, f, f, gocv.InterpolationLinear)
}

func ToGray(img gocv.Mat) gocv.Mat {
	gray := gocv.NewMat()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	return gray
}

func FindEdge(img gocv.Mat) gocv.Mat {
	edge := gocv.NewMat()
	gocv.Canny(img, &edge, 50, 150)
	return edge
}

func FindContours(img gocv.Mat) (gocv.PointsVector, gocv.Mat) {
	hierarchy := gocv.NewMat()
	contours := gocv.FindContoursWithParams(img, &hierarchy, gocv.RetrievalList, gocv.ChainApproxSimple)
	return contours, hierarchy
}

func FilterContours(contours gocv.PointsVector, min float64) gocv.PointsVector {
	largeContours := gocv.NewPointsVector()
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		if gocv.ContourArea(contour) > min {
			largeContours.Append(contour)
		}
	}
	return largeContours
}

func GetConvexes(contours gocv.PointsVector) []gocv.Mat {
	convexes := make([]gocv.Mat, contours.Size())
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)

		convex := gocv.NewMat()
		gocv.ConvexHull(contour, &convex, true, false) // convexはcontourの部分集合のインデックス
		convexes[i] = convex
	}
	return convexes
}

func GetPolygons(contours gocv.PointsVector, convexes []gocv.Mat) gocv.PointsVector {
	polies := gocv.NewPointsVector()
	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
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
	return polies
}

func SelectNearestSquareIndex(polies gocv.PointsVector) (int, float64) {
	selectedIndex := -1 // 選ばれたindex
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
	return selectedIndex, pMinError
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
