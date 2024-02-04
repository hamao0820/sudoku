package solver

import (
	"testing"
)

func TestFromArray(t *testing.T) {
	/*
		000200100
		145090000
		060800005
		000400000
		401602000
		002180039
		020000607
		000008050
		000500023
	*/
	b_ := [9][9]int{
		{0, 0, 0, 2, 0, 0, 1, 0, 0},
		{1, 4, 5, 0, 9, 0, 0, 0, 0},
		{0, 6, 0, 8, 0, 0, 0, 0, 5},
		{0, 0, 0, 4, 0, 0, 0, 0, 0},
		{4, 0, 1, 6, 0, 2, 0, 0, 0},
		{0, 0, 2, 1, 8, 0, 0, 3, 9},
		{0, 2, 0, 0, 0, 0, 6, 0, 7},
		{0, 0, 0, 0, 0, 8, 0, 5, 0},
		{0, 0, 0, 5, 0, 0, 0, 2, 3},
	}
	b1 := NewBoard()
	e := b1.FromArray(b_)
	if e != nil {
		t.Fatal(e)
	}
	b2 := NewBoard()
	b2[0][3] = 2
	b2[0][6] = 1
	b2[1][0] = 1
	b2[1][1] = 4
	b2[1][2] = 5
	b2[1][4] = 9
	b2[2][1] = 6
	b2[2][3] = 8
	b2[2][8] = 5
	b2[3][3] = 4
	b2[4][0] = 4
	b2[4][2] = 1
	b2[4][3] = 6
	b2[4][5] = 2
	b2[5][2] = 2
	b2[5][3] = 1
	b2[5][4] = 8
	b2[5][7] = 3
	b2[5][8] = 9
	b2[6][1] = 2
	b2[6][6] = 6
	b2[6][8] = 7
	b2[7][5] = 8
	b2[7][7] = 5
	b2[8][3] = 5
	b2[8][7] = 2
	b2[8][8] = 3

	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b1[y][x] != b2[y][x] {
				t.Errorf("Cells[%d][%d] = %d, want %d", y, x, b1[y][x], b2[y][x])
			}
		}
	}
}

func TestFromFile(t *testing.T) {
	/*
		000200100
		145090000
		060800005
		000400000
		401602000
		002180039
		020000607
		000008050
		000500023
	*/
	b1 := NewBoard()
	e := b1.FromFile("test/easy1.txt")
	if e != nil {
		t.Fatal(e)
	}
	b2 := NewBoard()
	b2[0][3] = 2
	b2[0][6] = 1
	b2[1][0] = 1
	b2[1][1] = 4
	b2[1][2] = 5
	b2[1][4] = 9
	b2[2][1] = 6
	b2[2][3] = 8
	b2[2][8] = 5
	b2[3][3] = 4
	b2[4][0] = 4
	b2[4][2] = 1
	b2[4][3] = 6
	b2[4][5] = 2
	b2[5][2] = 2
	b2[5][3] = 1
	b2[5][4] = 8
	b2[5][7] = 3
	b2[5][8] = 9
	b2[6][1] = 2
	b2[6][6] = 6
	b2[6][8] = 7
	b2[7][5] = 8
	b2[7][7] = 5
	b2[8][3] = 5
	b2[8][7] = 2
	b2[8][8] = 3
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b1[y][x] != b2[y][x] {
				t.Errorf("Cells[%d][%d] = %d, want %d", y, x, b1[y][x], b2[y][x])
			}
		}
	}
}

func TestVerify(t *testing.T) {
	{
		b := NewBoard()
		b[0][0] = 1
		b[0][1] = 2
		b[0][2] = 3
		b[0][3] = 4
		b[0][4] = 5
		b[0][5] = 6
		b[0][6] = 7
		b[0][7] = 8
		b[0][8] = 9
		if !b.Verify() {
			t.Error("Verify() = false, want true")
		}
		b[0][1] = 1
		if b.Verify() {
			t.Error("Verify() = true, want false")
		}
	}
	{
		b := NewBoard()
		b[0][0] = 1
		b[1][0] = 2
		b[2][0] = 3
		b[3][0] = 4
		b[4][0] = 5
		b[5][0] = 6
		b[6][0] = 7
		b[7][0] = 8
		b[8][0] = 9
		if !b.Verify() {
			t.Error("Verify() = false, want true")
		}
		b[1][0] = 1
		if b.Verify() {
			t.Error("Verify() = true, want false")
		}
	}
	{
		b := NewBoard()
		b[0][0] = 1
		b[0][1] = 2
		b[0][2] = 3
		b[1][0] = 4
		b[1][1] = 5
		b[1][2] = 6
		b[2][0] = 7
		b[2][1] = 8
		b[2][2] = 9
		if !b.Verify() {
			t.Error("Verify() = false, want true")
		}
		b[1][1] = 1
		if b.Verify() {
			t.Error("Verify() = true, want false")
		}
	}
}

func TestIsLegal(t *testing.T) {
	{
		b := NewBoard()
		b[0][0] = 1
		b[0][1] = 2
		b[0][2] = 3
		b[0][3] = 4
		b[0][4] = 5
		b[0][5] = 6
		b[0][6] = 7
		b[0][7] = 8
		if !b.IsLegal(0, 8, 9) {
			t.Error("IsLegal(0, 8, 9) = false, want true")
		}
		if b.IsLegal(0, 8, 1) {
			t.Error("IsLegal(0, 8, 1) = true, want false")
		}
	}
	{
		b := NewBoard()
		b[0][0] = 1
		b[1][0] = 2
		b[2][0] = 3
		b[3][0] = 4
		b[4][0] = 5
		b[5][0] = 6
		b[6][0] = 7
		b[7][0] = 8
		if !b.IsLegal(8, 0, 9) {
			t.Error("IsLegal(8, 0, 9) = false, want true")
		}
		if b.IsLegal(8, 0, 1) {
			t.Error("IsLegal(8, 0, 1) = true, want false")
		}
	}
	{
		b := NewBoard()
		b[0][0] = 1
		b[0][1] = 2
		b[0][2] = 3
		b[1][0] = 4
		b[1][1] = 5
		b[1][2] = 6
		b[2][0] = 7
		b[2][1] = 8
		if !b.IsLegal(2, 2, 9) {
			t.Error("IsLegal(2, 2, 9) = false, want true")
		}
		if b.IsLegal(2, 2, 1) {
			t.Error("IsLegal(2, 2, 1) = true, want false")
		}
	}
}

func TestIsFull(t *testing.T) {
	{
		b := NewBoard()
		if b.IsFull() {
			t.Error("IsFull() = true, want false")
		}
	}
	{
		b := NewBoard()
		b[0][0] = 1
		if b.IsFull() {
			t.Error("IsFull() = true, want false")
		}
	}
	{
		b := NewBoard()
		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				b[x][y] = 1
			}
		}
		if !b.IsFull() {
			t.Error("IsFull() = false, want true")
		}
	}
}

func TestIsSolved(t *testing.T) {
	{
		b := NewBoard()
		if b.IsSolved() {
			t.Error("IsSolved() = true, want false")
		}
	}
	{
		b := NewBoard()
		b[0][0] = 1
		if b.IsSolved() {
			t.Error("IsSolved() = true, want false")
		}
	}
	{
		b := NewBoard()
		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				b[x][y] = 1
			}
		}
		if b.IsSolved() {
			t.Error("IsSolved() = true, want false")
		}
	}
	{
		b := NewBoard()
		b_ := [9][9]int{
			{1, 7, 4, 2, 3, 8, 6, 9, 5},
			{6, 3, 9, 7, 4, 5, 2, 8, 1},
			{2, 8, 5, 1, 9, 6, 7, 3, 4},
			{4, 1, 7, 8, 2, 3, 9, 5, 6},
			{5, 9, 2, 6, 1, 4, 8, 7, 3},
			{8, 6, 3, 5, 7, 9, 1, 4, 2},
			{7, 4, 1, 3, 8, 2, 5, 6, 9},
			{3, 5, 8, 9, 6, 1, 4, 2, 7},
			{9, 2, 6, 4, 5, 7, 3, 1, 8},
		}
		for y := 0; y < 9; y++ {
			for x := 0; x < 9; x++ {
				b[x][y] = b_[y][x]
			}
		}
		if !b.IsSolved() {
			t.Error("IsSolved() = false, want true")
		}
	}
}
