package detect

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func FitSize(img *gocv.Mat, h, w int) {
	size := img.Size()
	f := math.Min(float64(h)/float64(size[0]), float64(w)/float64(size[1]))
	gocv.Resize(*img, img, image.Point{}, f, f, gocv.InterpolationLinear)
}

func Edge(img gocv.Mat) gocv.Mat {
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
	edgeImage := gocv.NewMat()
	gocv.Canny(gray, &edgeImage, 50, 150)
	return edgeImage
}

// def line(img, show=True, threshold=80, minLineLength=50, maxLineGap=5):
//     edges = edge(img, False)
//     lines = cv2.HoughLinesP(edges, 1, np.pi/180, threshold, 200, minLineLength, maxLineGap)
//     return lines

func Line(img gocv.Mat, threshold int) gocv.Mat {
	edgeImage := Edge(img)
	defer edgeImage.Close()
	lines := gocv.NewMat()
	gocv.HoughLinesP(edgeImage, &lines, 1, math.Pi/180, threshold)
	return lines
}

// def contours(img, show=True):
//     edges = edge(img, False)
//     contours = cv2.findContours(edges, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)[1]
//     blank = np.zeros(img.shape, np.uint8)
//     min_area = img.shape[0] * img.shape[1] * 0.2 # 画像の何割占めるか
//     large_contours = [c for c in contours if cv2.contourArea(c) > min_area]
//     cv2.drawContours(blank, large_contours, -1, (0,255,0), 1)
//     return large_contours

func Contours(img gocv.Mat) (gocv.PointsVector, gocv.Mat) {
	edges := Edge(img)
	contours := gocv.FindContours(edges, gocv.RetrievalList, gocv.ChainApproxSimple)
	blank := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	min_area := float64(img.Rows()*img.Cols()) * 0.2
	large_contours := gocv.NewPointsVector()
	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)
		if gocv.ContourArea(cnt) > min_area {
			large_contours.Append(cnt)
		}
	}
	gocv.DrawContours(&blank, large_contours, -1, color.RGBA{0, 255, 0, 255}, 1)
	return large_contours, blank
}

// def convex(img, show=True):
//     blank = np.copy(img)
//     convexes = []
//     for cnt in contours(img, False):
//         convex = cv2.convexHull(cnt)
//         cv2.drawContours(blank, [convex], -1, (0,255,0), 2)
//         convexes.append(convex)
//     return convexes

func Convex(img gocv.Mat) gocv.PointsVector {
	blank := gocv.NewMat()
	defer blank.Close()
	img.CopyTo(&blank)
	contours, _ := Contours(img)
	convexes := gocv.NewPointsVector()
	for i := 0; i < contours.Size(); i++ {
		cnt := contours.At(i)
		convex := gocv.NewMat()
		defer convex.Close()
		gocv.ConvexHull(cnt, &convex, false, false)
		convexes.Append(cnt)
	}
	return convexes
}

// def convex_poly(img, show=True):
//     cnts = convex(img, False)
//     blank = np.copy(img)
//     polies = []
//     for cnt in cnts:
//         arclen = cv2.arcLength(cnt, True)
//         poly = cv2.approxPolyDP(cnt, 0.02*arclen, True)
//         cv2.drawContours(blank, [poly], -1, (0,255,0), 2)
//         polies.append(poly)
//     return [poly[:, 0, :] for poly in polies]

func convex_poly(img gocv.Mat) gocv.PointsVector {
	cnts := Convex(img)
	blank := gocv.NewMat()
	defer blank.Close()
	img.CopyTo(&blank)
	polies := gocv.NewPointsVector()
	for i := 0; i < cnts.Size(); i++ {
		cnt := cnts.At(i)
		arclen := gocv.ArcLength(cnt, true)
		poly := gocv.ApproxPolyDP(cnt, 0.02*arclen, true)
		polies.Append(poly)
	}
	gocv.DrawContours(&blank, polies, -1, color.RGBA{0, 255, 0, 255}, 2)
	return polies
}

// def select_corners(img, polies):
//     p_selected = []
//     p_scores = []
//     for poly in polies:
//         choices = np.array(list(itertools.combinations(poly, 4)))
//         scores = []
//         # 正方形に近いものを選ぶ
//         for c in choices:
//             line_lens = [np.linalg.norm(c[(i + 1) % 4] - c[i]) for i in range(4)]
//             base = cv2.contourArea(c) ** 0.5
//             score = sum([abs(1 - l/base) ** 2 for l in line_lens])
//             scores.append(score)
//         idx = np.argmin(scores)
//         p_selected.append(choices[idx])
//         p_scores.append(scores[idx])
//     return p_selected[np.argmin(p_scores)]

func SelectCorners(img gocv.Mat, polies gocv.PointsVector) gocv.PointVector {
	p_selected := []gocv.PointVector{}
	p_scores := []float64{}
	for i := 0; i < polies.Size(); i++ {
		poly := polies.At(i)
		choices := [][]image.Point{}
		for i := 0; i < 4; i++ {
			choices = append(choices, []image.Point{poly.At(i), poly.At((i + 1) % 4), poly.At((i + 2) % 4), poly.At((i + 3) % 4)})
		}
		scores := []float64{}
		for _, c := range choices {
			line_lens := []float64{}
			for i := 0; i < 4; i++ {
				a := c[(i+1)%4].Sub(c[i])
				line_lens = append(line_lens, math.Sqrt(float64(a.X*a.X+a.Y*a.Y)))
			}
			base := math.Sqrt(gocv.ContourArea(gocv.NewPointVectorFromPoints(c)))
			score := 0.0
			for _, l := range line_lens {
				score += math.Pow(math.Abs(1-l/base), 2)
			}
			scores = append(scores, score)
		}
		idx := 0
		min_score := scores[0]
		for i, score := range scores {
			if score < min_score {
				idx = i
				min_score = score
			}
		}
		p_selected = append(p_selected, gocv.NewPointVectorFromPoints(choices[idx]))
		p_scores = append(p_scores, scores[idx])
	}
	min_score := p_scores[0]
	min_idx := 0
	for i, score := range p_scores {
		if score < min_score {
			min_score = score
			min_idx = i
		}
	}
	return p_selected[min_idx]
}

// def gen_score_mat():
//     half_a = np.fromfunction(lambda i, j: ((10 - i) ** 2) / 100.0, (10, 20), dtype=np.float32)
//     half_b = np.rot90(half_a, 2)
//     cell_a = np.r_[half_a, half_b]
//     cell_b = np.rot90(cell_a)
//     cell = np.maximum(cell_a, cell_b)
//     return np.tile(cell, (9, 9))

func GenScoreMat() gocv.Mat {
	half_a := gocv.NewMatWithSize(10, 20, gocv.MatTypeCV32F)
	for i := 0; i < 10; i++ {
		for j := 0; j < 20; j++ {
			half_a.SetFloatAt(i, j, float32((10-i)*(10-i))/100.0)
		}
	}
	half_b := gocv.NewMat()
	defer half_b.Close()
	gocv.Rotate(half_a, &half_b, gocv.Rotate180Clockwise)
	cell_a := gocv.NewMat()
	defer cell_a.Close()
	gocv.Vconcat(half_a, half_b, &cell_a)
	cell_b := gocv.NewMat()
	defer cell_b.Close()
	gocv.Rotate(cell_a, &cell_b, gocv.Rotate90Clockwise)
	cell := gocv.NewMat()
	defer cell.Close()
	gocv.Max(cell_a, cell_b, &cell)
	score_mat := gocv.NewMatWithSize(cell.Rows()*9, cell.Cols()*9, gocv.MatTypeCV32F)
	defer score_mat.Close()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			roi := cell.Region(image.Rect(0, 0, cell.Cols(), cell.Rows()))
			for k := 0; k < cell.Rows(); k++ {
				for l := 0; l < cell.Cols(); l++ {
					score_mat.SetFloatAt(i*cell.Rows()+k, j*cell.Cols()+l, roi.GetFloatAt(k, l))
				}
			}
		}
	}
	return score_mat
}

// SCALE = 0.7
// def get_get_fit_score(img, x):
//     # 入力リサイズ
//     img = cv2.resize(img, (int(img.shape[1] * SCALE), int(img.shape[0] * SCALE)), interpolation=cv2.INTER_AREA)
//     img_size = (img.shape[0] * img.shape[1]) ** 0.5
//     x = np.int32(x * SCALE)

//     # 線分化
//     poly_length = cv2.arcLength(x, True)
//     lines = line(img, False, int(poly_length / 12), int(poly_length / 200))
//     line_mat = np.zeros(img.shape, np.uint8)
//     for x1, y1, x2, y2 in lines[:, 0]:
//         cv2.line(line_mat, (x1, y1), (x2, y2), 255, 1)
//     line_mat = line_mat[:, :, 0]

//     # 矩形の外をマスクアウト
//     img_size = (img.shape[0] * img.shape[1]) ** 0.5
//     mask = np.zeros(line_mat.shape, np.uint8)
//     cv2.fillConvexPoly(mask, x, 1)
//     kernel = np.ones((int(img_size / 10), int(img_size / 10)), np.uint8)
//     mask = cv2.erode(mask, kernel, iterations=1)
//     line_mat[np.where(mask == 0)] = 0

//     # スコア
//     score_mat = gen_score_mat()

//     def get_fit_score(x):
//         img_pnts = np.float32(x).reshape(4, 2)
//         img_pnts *= SCALE
//         score_size = score_mat.shape[0]
//         score_pnts = np.float32([[0, 0], [0, score_size], [score_size, score_size], [score_size, 0]])

//         transform = cv2.getPerspectiveTransform(score_pnts, img_pnts)
//         score_t = cv2.warpPerspective(score_mat, transform, (img.shape[1], img.shape[0]))

//         res = line_mat * score_t
//         return -np.average(res[np.where(res > 255 * 0.1)])

//     return get_fit_score

// 	poly_length := gocv.ArcLength(x_int32, true)
// 	lines := line(img_resized, int(poly_length/12))
// 	line_mat := gocv.NewMatWithSize(img_resized.Rows(), img_resized.Cols(), gocv.MatTypeCV8U)
// 	defer line_mat.Close()
// 	for i := 0; i < lines.Rows(); i++ {
// 		line := lines.GetVeciAt(i, 0)
// 		gocv.Line(&line_mat, image.Pt(line.Val1, line.Val2), image.Pt(line.Val3, line.Val4), color.RGBA{255, 255, 255, 255}, 1)
// 	}

// }

// def convex_poly_fitted(img, show=True):
//     polies = convex_poly(img, False)
//     poly = select_corners(img, polies)
//     x0 = poly.flatten()
//     get_fit_score = get_get_fit_score(img, poly)
//     ret = basinhopping(get_fit_score, x0, T=0.1, niter=250, stepsize=3)
//     return ret.x.reshape(4, 2), ret.fun

// func ConvexPolyFitted(img gocv.Mat) (gocv.PointVector, float64) {
// 	polies := convex_poly(img)
// 	poly := SelectCorners(img, polies)
// 	x0 := poly.ToPoints()
// 	get_fit_score := get_get_fit_score(img, poly)
// 	ret := basinhopping(get_fit_score, x0, 0.1, 250, 3)
// 	return gocv.NewPointVectorFromPoints(ret.X), ret.Fun
// }

// func convex_poly_fitted(img gocv.Mat) (gocv.PointVector, float64) {
// 	polies := convex_poly(img)
// 	poly := select_corners(img, polies)
// 	x0 := poly.ToPoints()
// 	get_fit_score := get_get_fit_score(img, poly)
// 	ret := basinhopping(get_fit_score, x0, 0.1, 250, 3)
// 	return gocv.NewPointVectorFromPoints(ret.X), ret.Fun
// }

// def normalize_corners(v):
//     rads = []
//     for i in range(4):
//         a = v[(i + 1) % 4] - v[i]
//         a = a / np.linalg.norm(a)
//         cosv = np.dot(a, np.array([1, 0]))
//         rads.append(math.acos(cosv))
//     left_top = np.argmin(rads)
//     return np.roll(v, 4 - left_top, axis=0)
