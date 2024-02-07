package ocr

import (
	"fmt"
	"os"
	"time"

	"github.com/hamao0820/sudoku/detect"
	"gocv.io/x/gocv"
)

func CollectCell() {
	webcam, _ := gocv.OpenVideoCapture(0)
	webcam.Set(gocv.VideoCaptureFPS, 30)
	defer webcam.Close()
	window := gocv.NewWindow("Test")
	defer window.Close()
	img := gocv.NewMat()
	defer img.Close()

	// iをcacheする
	f, err := os.Open("ocr/cache.txt")
	if err != nil {
		panic(err)
	}
	var i int
	fmt.Fscanf(f, "%d", &i)
	f.Close()
	for {
		webcam.Read(&img)

		square, err := detect.GetSquare(img)
		if err == nil {
			cells := detect.GetCells(square)
			for y := 0; y < 9; y++ {
				for x := 0; x < 9; x++ {
					gocv.IMWrite(fmt.Sprintf("ocr/data/raw/%4d.png", i), cells[y][x])
					i++
				}
			}

			f, err := os.Create("ocr/cache.txt")
			if err != nil {
				panic(err)
			}
			_, err = f.Write([]byte(fmt.Sprintf("%d", i)))
			if err != nil {
				panic(err)
			}
			fmt.Println("Saved", i)
			time.Sleep(1 * time.Second)
		}

		window.IMShow(img)
		window.WaitKey(1)

		time.Sleep(300 * time.Millisecond)
	}
}
