package detect

import (
	"testing"

	"gocv.io/x/gocv"
)

func TestFitSize(t *testing.T) {
	{
		img := gocv.NewMatWithSize(300, 300, gocv.MatTypeCV8U)
		defer img.Close()

		FitSize(&img, 100, 100)
		if img.Size()[0] != 100 || img.Size()[1] != 100 {
			t.Error("fit_size failed")
		}
	}
	{
		img := gocv.NewMatWithSize(300, 300, gocv.MatTypeCV8U)
		defer img.Close()

		FitSize(&img, 200, 100)
		if img.Size()[0] != 100 || img.Size()[1] != 100 {
			t.Error("fit_size failed")
		}
	}
	{
		img := gocv.NewMatWithSize(300, 500, gocv.MatTypeCV8U)
		defer img.Close()

		FitSize(&img, 100, 100)
		if img.Size()[0] != 60 || img.Size()[1] != 100 {
			t.Error("fit_size failed")
		}
	}
	{
		img := gocv.NewMatWithSize(100, 200, gocv.MatTypeCV8U)
		defer img.Close()

		FitSize(&img, 200, 300)
		if img.Size()[0] != 150 || img.Size()[1] != 300 {
			t.Error("fit_size failed")
		}
	}
}
