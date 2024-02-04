package solver

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Board [9][9]int

func NewBoard() *Board {
	board := Board{}
	return &board
}

func FromFile(filename string) (*Board, error) {
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
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	b := NewBoard()
	for y := 0; y < 9; y++ {
		scanner.Scan()
		line := scanner.Text()
		if len(line) != 9 {
			return nil, errors.New("invalid file format")
		}
		for x := 0; x < 9; x++ {
			if line[x] == '0' {
				b[y][x] = 0
			} else if '1' <= line[x] && line[x] <= '9' {
				b[y][x] = int(line[x] - 48)
			} else {
				return nil, errors.New("invalid file format")
			}
		}
	}
	return b, nil
}

func (b *Board) IsEmpty(y, x int) bool {
	return b[y][x] == 0
}

func (b *Board) IsFull() bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b[y][x] == 0 {
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
		if i != x && b[y][i] == value {
			return false
		}
		if i != y && b[i][x] == value {
			return false
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if (i != y%3 || j != x%3) && b[y/3*3+i][x/3*3+j] == value {
				return false
			}
		}
	}
	return true
}

func (b *Board) Verify() bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b[y][x] == 0 {
				continue
			}
			if !b.IsLegal(y, x, b[y][x]) {
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
			if b[y][x] != 0 {
				str += fmt.Sprint(b[y][x])
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
