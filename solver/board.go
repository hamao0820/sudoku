package solver

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Cell struct {
	Value   int
	IsFixed bool
}

type Board struct {
	Cells [9][9]Cell
}

func NewBoard() *Board {
	board := Board{}
	return &board
}

func (b *Board) FromArray(array [9][9]int) error {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			b.Cells[y][x].Value = array[y][x]
			if array[y][x] != 0 {
				b.Cells[y][x].IsFixed = true
			}
		}
	}
	return nil
}

func (b *Board) FromFile(filename string) error {
	/*
		入力ファイル例:
		000200100
		145090000
		060800005
		000400000
		401602000
		002180039
		020000607
		000008050
		000500023
		各行は9文字で、空白は0、数字は1-9である。
	*/
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for y := 0; y < 9; y++ {
		scanner.Scan()
		line := scanner.Text()
		if len(line) != 9 {
			return errors.New("invalid file format")
		}
		for x := 0; x < 9; x++ {
			if line[x] == '0' {
				b.Cells[y][x].Value = 0
			} else if '1' <= line[x] && line[x] <= '9' {
				b.Cells[y][x].Value = int(line[x] - 48)
				b.Cells[y][x].IsFixed = true
			} else {
				return errors.New("invalid file format")
			}
		}
	}
	return nil
}

func (b *Board) IsEmpty(y, x int) bool {
	return b.Cells[y][x].Value == 0
}

func (b *Board) IsFull() bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b.Cells[y][x].Value == 0 {
				return false
			}
		}
	}
	return true
}

func (b *Board) IsSolved() bool {
	if !b.IsFull() {
		return false
	}
	return b.Verify()
}

func (b *Board) IsLegal(y, x, value int) bool {
	for i := 0; i < 9; i++ {
		if i != x && b.Cells[y][i].Value == value {
			return false
		}
		if i != y && b.Cells[i][x].Value == value {
			return false
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if (i != y%3 || j != x%3) && b.Cells[y/3*3+i][x/3*3+j].Value == value {
				return false
			}
		}
	}
	return true
}

func (b *Board) Verify() bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b.Cells[y][x].Value == 0 {
				continue
			}
			if !b.IsLegal(y, x, b.Cells[y][x].Value) {
				return false
			}
		}
	}
	return true
}

func (b *Board) String() string {
	str := ""
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b.Cells[y][x].Value != 0 {
				str += fmt.Sprint(b.Cells[y][x].Value)
			} else {
				str += "･"
			}
			if x == 2 || x == 5 {
				str += " "
			}
		}
		str += "\n"
		if y == 2 || y == 5 {
			str += "\n"
		}
	}
	return str
}
