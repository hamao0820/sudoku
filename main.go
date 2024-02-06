package main

import (
	"image"
	"image/color"

	"github.com/hamao0820/sudoku/detect"
	"gocv.io/x/gocv"
)

func main() {
	src := gocv.IMRead("sample.png", gocv.IMReadColor)
	defer src.Close()

	detect.FitSize(&src, 500, 500)

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(src, &gray, gocv.ColorBGRToGray)

	binary := gocv.NewMat()
	defer binary.Close()
	gocv.Threshold(gray, &binary, 100, 255, gocv.ThresholdToZeroInv)
	gocv.BitwiseNot(binary, &binary)
	gocv.Threshold(gray, &binary, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	edge := gocv.NewMat()
	defer edge.Close()
	gocv.Canny(binary, &edge, 50, 150)

	hierarchy := gocv.NewMat()
	defer hierarchy.Close()
	contours := gocv.FindContoursWithParams(edge, &hierarchy, gocv.RetrievalList, gocv.ChainApproxSimple)

	warp := gocv.NewMat()
	defer warp.Close()

	min_area := float64(src.Rows()*src.Cols()) * 0.2
	max_level := 0
	max_area := 0.0
	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		if area < min_area {
			continue
		}

		arcLen := gocv.ArcLength(contours.At(i), true)
		approx := gocv.ApproxPolyDP(contours.At(i), 0.02*arcLen, true)
		if approx.Size() != 4 {
			continue
		}

		if area > max_area {
			max_area = area
		}
	}

	trimmed := gocv.NewMat()
	defer trimmed.Close()

	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		if area < max_area {
			continue
		}

		arcLen := gocv.ArcLength(contours.At(i), true)
		approx := gocv.ApproxPolyDP(contours.At(i), 0.02*arcLen, true)
		if approx.Size() != 4 {
			continue
		}

		dst := gocv.NewPointVector()
		// 300x300
		// 反転を修正
		dst.Append(image.Pt(0, 0))
		dst.Append(image.Pt(300, 0))
		dst.Append(image.Pt(300, 300))
		dst.Append(image.Pt(0, 300))

		defer dst.Close()
		// gocv.DrawContoursWithParams(&src, contours, i, color.RGBA{255, 0, 0, 255}, 1, gocv.LineAA, hierarchy, max_level, image.Pt(0, 0))
		gocv.WarpPerspective(src, &trimmed, gocv.GetPerspectiveTransform(approx, dst), image.Pt(300, 300))
		// 左右反転
		gocv.Flip(trimmed, &trimmed, 1)

		perspective := gocv.GetPerspectiveTransform(approx, dst)
		gocv.WarpPerspective(src, &warp, perspective, image.Pt(src.Cols(), src.Rows()))

		gocv.DrawContoursWithParams(&src, contours, i, color.RGBA{255, 0, 0, 255}, 1, gocv.LineAA, hierarchy, max_level, image.Pt(0, 0))
		break
	}

	win := gocv.NewWindow("sudoku")
	defer win.Close()
	win.IMShow(src)
	win.WaitKey(0)

	// edges := gocv.NewMat()
	// defer edges.Close()
	// gocv.Canny(gray, &edges, 50, 150)

	// lines := gocv.NewMat()
	// defer lines.Close()
	// gocv.HoughLinesPWithParams(edges, &lines, 1, math.Pi/180, 80, 50, 5)

	// contours := gocv.FindContours(edges, gocv.RetrievalList, gocv.ChainApproxSimple)
	// minArea := float64(src.Rows()*src.Cols()) * 0.2
	// lergeContours := gocv.NewPointsVector()
	// for i := 0; i < contours.Size(); i++ {
	// 	cnt := contours.At(i)
	// 	if gocv.ContourArea(cnt) > minArea {
	// 		lergeContours.Append(cnt)
	// 	}
	// }

	// hulls := []gocv.Mat{}
	// for i := 0; i < lergeContours.Size(); i++ {
	// 	cnt := lergeContours.At(i)
	// 	hull := gocv.NewMat()
	// 	gocv.ConvexHull(cnt, &hull, true, false)
	// 	hulls = append(hulls, hull)
	// }

	// blank := gocv.NewMatWithSize(src.Rows(), src.Cols(), gocv.MatTypeCV8U)
	// polies := gocv.NewPointsVector()
	// for _, hull := range hulls {
	// 	pv := gocv.NewPointVectorFromMat(hull)
	// 	arcLen := gocv.ArcLength(pv, true)
	// 	poly := gocv.ApproxPolyDP(pv, 0.02*arcLen, true)
	// 	gocv.DrawContours(&blank, gocv.NewPointsVectorFromPoints([][]image.Point{poly.ToPoints()}), -1, color.RGBA{0, 255, 0, 255}, 1)
	// 	polies.Append(poly)
	// }

	// win := gocv.NewWindow("sudoku")
	// defer win.Close()

	// win.IMShow(blank)
	// win.WaitKey(0)
}
