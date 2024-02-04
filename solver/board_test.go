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
	b2.Cells[0][3].Value = 2
	b2.Cells[0][6].Value = 1
	b2.Cells[1][0].Value = 1
	b2.Cells[1][1].Value = 4
	b2.Cells[1][2].Value = 5
	b2.Cells[1][4].Value = 9
	b2.Cells[2][1].Value = 6
	b2.Cells[2][3].Value = 8
	b2.Cells[2][8].Value = 5
	b2.Cells[3][3].Value = 4
	b2.Cells[4][0].Value = 4
	b2.Cells[4][2].Value = 1
	b2.Cells[4][3].Value = 6
	b2.Cells[4][5].Value = 2
	b2.Cells[5][2].Value = 2
	b2.Cells[5][3].Value = 1
	b2.Cells[5][4].Value = 8
	b2.Cells[5][7].Value = 3
	b2.Cells[5][8].Value = 9
	b2.Cells[6][1].Value = 2
	b2.Cells[6][6].Value = 6
	b2.Cells[6][8].Value = 7
	b2.Cells[7][5].Value = 8
	b2.Cells[7][7].Value = 5
	b2.Cells[8][3].Value = 5
	b2.Cells[8][7].Value = 2
	b2.Cells[8][8].Value = 3

	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b1.Cells[y][x].Value != b2.Cells[y][x].Value {
				t.Errorf("Cells[%d][%d].Value = %d, want %d", y, x, b1.Cells[y][x].Value, b2.Cells[y][x].Value)
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
	b2.Cells[0][3].Value = 2
	b2.Cells[0][6].Value = 1
	b2.Cells[1][0].Value = 1
	b2.Cells[1][1].Value = 4
	b2.Cells[1][2].Value = 5
	b2.Cells[1][4].Value = 9
	b2.Cells[2][1].Value = 6
	b2.Cells[2][3].Value = 8
	b2.Cells[2][8].Value = 5
	b2.Cells[3][3].Value = 4
	b2.Cells[4][0].Value = 4
	b2.Cells[4][2].Value = 1
	b2.Cells[4][3].Value = 6
	b2.Cells[4][5].Value = 2
	b2.Cells[5][2].Value = 2
	b2.Cells[5][3].Value = 1
	b2.Cells[5][4].Value = 8
	b2.Cells[5][7].Value = 3
	b2.Cells[5][8].Value = 9
	b2.Cells[6][1].Value = 2
	b2.Cells[6][6].Value = 6
	b2.Cells[6][8].Value = 7
	b2.Cells[7][5].Value = 8
	b2.Cells[7][7].Value = 5
	b2.Cells[8][3].Value = 5
	b2.Cells[8][7].Value = 2
	b2.Cells[8][8].Value = 3
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b1.Cells[y][x].Value != b2.Cells[y][x].Value {
				t.Errorf("Cells[%d][%d].Value = %d, want %d", y, x, b1.Cells[y][x].Value, b2.Cells[y][x].Value)
			}
		}
	}
}

func TestVerify(t *testing.T) {
	{
		b := NewBoard()
		b.Cells[0][0].Value = 1
		b.Cells[0][1].Value = 2
		b.Cells[0][2].Value = 3
		b.Cells[0][3].Value = 4
		b.Cells[0][4].Value = 5
		b.Cells[0][5].Value = 6
		b.Cells[0][6].Value = 7
		b.Cells[0][7].Value = 8
		b.Cells[0][8].Value = 9
		if !b.Verify() {
			t.Error("Verify() = false, want true")
		}
		b.Cells[0][1].Value = 1
		if b.Verify() {
			t.Error("Verify() = true, want false")
		}
	}
	{
		b := NewBoard()
		b.Cells[0][0].Value = 1
		b.Cells[1][0].Value = 2
		b.Cells[2][0].Value = 3
		b.Cells[3][0].Value = 4
		b.Cells[4][0].Value = 5
		b.Cells[5][0].Value = 6
		b.Cells[6][0].Value = 7
		b.Cells[7][0].Value = 8
		b.Cells[8][0].Value = 9
		if !b.Verify() {
			t.Error("Verify() = false, want true")
		}
		b.Cells[1][0].Value = 1
		if b.Verify() {
			t.Error("Verify() = true, want false")
		}
	}
	{
		b := NewBoard()
		b.Cells[0][0].Value = 1
		b.Cells[0][1].Value = 2
		b.Cells[0][2].Value = 3
		b.Cells[1][0].Value = 4
		b.Cells[1][1].Value = 5
		b.Cells[1][2].Value = 6
		b.Cells[2][0].Value = 7
		b.Cells[2][1].Value = 8
		b.Cells[2][2].Value = 9
		if !b.Verify() {
			t.Error("Verify() = false, want true")
		}
		b.Cells[1][1].Value = 1
		if b.Verify() {
			t.Error("Verify() = true, want false")
		}
	}
}

func TestIsLegal(t *testing.T) {
	{
		b := NewBoard()
		b.Cells[0][0].Value = 1
		b.Cells[0][1].Value = 2
		b.Cells[0][2].Value = 3
		b.Cells[0][3].Value = 4
		b.Cells[0][4].Value = 5
		b.Cells[0][5].Value = 6
		b.Cells[0][6].Value = 7
		b.Cells[0][7].Value = 8
		if !b.IsLegal(0, 8, 9) {
			t.Error("IsLegal(0, 8, 9) = false, want true")
		}
		if b.IsLegal(0, 8, 1) {
			t.Error("IsLegal(0, 8, 1) = true, want false")
		}
	}
	{
		b := NewBoard()
		b.Cells[0][0].Value = 1
		b.Cells[1][0].Value = 2
		b.Cells[2][0].Value = 3
		b.Cells[3][0].Value = 4
		b.Cells[4][0].Value = 5
		b.Cells[5][0].Value = 6
		b.Cells[6][0].Value = 7
		b.Cells[7][0].Value = 8
		if !b.IsLegal(8, 0, 9) {
			t.Error("IsLegal(8, 0, 9) = false, want true")
		}
		if b.IsLegal(8, 0, 1) {
			t.Error("IsLegal(8, 0, 1) = true, want false")
		}
	}
	{
		b := NewBoard()
		b.Cells[0][0].Value = 1
		b.Cells[0][1].Value = 2
		b.Cells[0][2].Value = 3
		b.Cells[1][0].Value = 4
		b.Cells[1][1].Value = 5
		b.Cells[1][2].Value = 6
		b.Cells[2][0].Value = 7
		b.Cells[2][1].Value = 8
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
		b.Cells[0][0].Value = 1
		if b.IsFull() {
			t.Error("IsFull() = true, want false")
		}
	}
	{
		b := NewBoard()
		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				b.Cells[x][y].Value = 1
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
		b.Cells[0][0].Value = 1
		if b.IsSolved() {
			t.Error("IsSolved() = true, want false")
		}
	}
	{
		b := NewBoard()
		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				b.Cells[x][y].Value = 1
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
				b.Cells[x][y].Value = b_[y][x]
			}
		}
		if !b.IsSolved() {
			t.Error("IsSolved() = false, want true")
		}
	}
}
