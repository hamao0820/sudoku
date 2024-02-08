package ocr

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hamao0820/sudoku/detect"
	"gocv.io/x/gocv"
)

type SudokuID int

const (
	_ SudokuID = iota
	SudokuIDOne
	SudokuIDTwo
	SudokuIDThree
	SudokuIDFour
	SudokuIDFive
	SudokuIDSix
	SudokuIDSeven
	SudokuIDEight
	SudokuIDNine
	SudokuIDTen
)

func CollectSquareFromImage(id SudokuID, path string) {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading image from: ", path)
		return
	}
	defer img.Close()

	square, err := detect.GetSquare(img)
	if err != nil {
		fmt.Println(err)
		return
	}

	gocv.IMWrite(fmt.Sprintf("ocr/data/squares/%d/%d.png", id, rand.Int31()), square)
}

func CollectSquareFromCamera(id SudokuID) {
	fmt.Println("CollectSquareFromCamera")
	fmt.Println("id: ", id)

	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("Test")
	defer window.Close()
	img := gocv.NewMat()
	defer img.Close()

	i := 0
	frame := 0
	for {
		time.Sleep(100 * time.Millisecond)
		webcam.Read(&img)

		origin := gocv.NewMat()
		img.CopyTo(&origin)
		if drawSquare(&img) {
			frame++
			if frame >= 3 {
				square, err := detect.GetSquare(origin)
				if err != nil {
					continue
				}
				i++
				fmt.Printf("i: %02d\n", i)
				gocv.IMWrite(fmt.Sprintf("ocr/data/squares/%d/%d.png", id, rand.Int31()), square)
				frame = 0

				time.Sleep(50 * time.Millisecond)
			}
		} else {
			frame = 0
		}

		window.IMShow(img)
		window.WaitKey(1)
	}
}

func drawSquare(img *gocv.Mat) bool {
	detect.FitSize(img, 500, 500)

	gray := detect.ToGray(*img)
	defer gray.Close()

	edge := detect.FindEdge(gray)
	defer edge.Close()

	contours, _ := detect.FindContours(edge)
	defer contours.Close()

	min_area := float64(img.Rows()*img.Cols()) * 0.2
	largeContours := detect.FilterContours(contours, min_area)

	convexes := detect.GetConvexes(largeContours)

	polies := detect.GetPolygons(largeContours, convexes)
	defer polies.Close()

	// 正方形に近いものを選ぶ
	selectedIndex, _ := detect.SelectNearestSquareIndex(polies)
	if selectedIndex == -1 {
		return false
	}

	gocv.DrawContours(img, polies, selectedIndex, color.RGBA{255, 0, 0, 255}, 3)
	return true
}
