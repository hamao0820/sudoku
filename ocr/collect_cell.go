package ocr

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strconv"

	"gocv.io/x/gocv"
)

type Answer struct {
	Data Data `json:"data"`
}

type Data []struct {
	Id   SudokuID
	Cell Cell
}

type Cell [][]int

var (
	data = loadJSON()
)

func CollectCell() {
	paths, err := filepath.Glob("ocr/data/squares/*")
	if err != nil {
		panic(err)
	}
	counts := make(map[int]int)
	for _, path := range paths {
		if !isDir(path) {
			continue
		}
		dir := filepath.Base(path)
		id_, err := strconv.Atoi(dir)
		if err != nil {
			panic(err)
		}
		ansCells, err := selectCell(SudokuID(id_))
		if err != nil {
			fmt.Println(err)
			continue
		}
		images, err := filepath.Glob(filepath.Join(path, "*.png"))
		if err != nil {
			panic(err)
		}
		for _, image := range images {
			img := gocv.IMRead(image, gocv.IMReadGrayScale)
			if img.Empty() {
				fmt.Println("Error reading image from: ", image)
				continue
			}
			defer img.Close()

			cells := splitCells(img)
			for y := 0; y < 9; y++ {
				for x := 0; x < 9; x++ {
					cell := cells[y][x]
					if cell.Empty() {
						continue
					}
					ans := ansCells[y][x]
					gocv.IMWrite(fmt.Sprintf("ocr/data/cells/%d/%05d.png", ans, counts[ans]+1), cell)
					counts[ans]++
				}
			}
		}
	}
	for k, v := range counts {
		fmt.Printf("%d: %d\n", k, v)
	}
}

func splitCells(square gocv.Mat) [][]gocv.Mat {
	cells := make([][]gocv.Mat, 9)
	dx := float64(square.Cols()) / 9
	dy := float64(square.Rows()) / 9

	padding := 1.0
	for y := 0; y < 9; y++ {
		cells[y] = make([]gocv.Mat, 9)
		for x := 0; x < 9; x++ {
			cells[y][x] = square.Region(image.Rect(int(float64(x)*dx+padding), int(float64(y)*dy+padding), int(float64(x+1)*dx-padding), int(float64(y+1)*dy-padding)))
		}
	}
	return cells
}

func loadJSON() Data {
	f, err := os.Open("ocr/ans2.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ans Answer
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&ans); err != nil {
		panic(err)
	}

	return ans.Data
}

func selectCell(id SudokuID) (Cell, error) {
	for _, d := range data {
		if d.Id == id {
			return d.Cell, nil
		}
	}
	return [][]int{}, fmt.Errorf("not find")
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
