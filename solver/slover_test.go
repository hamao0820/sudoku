package solver

import (
	"testing"
)

func TestBacktrack(t *testing.T) {
	b := NewBoard()
	err := b.FromFile("test/easy1.txt")
	if err != nil {
		t.Fatal(err)
	}
	Backtrack(b)
	b2 := NewBoard()
	b2.FromFile("test/easy1_solution.txt")
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if b.Cells[y][x].Value != b2.Cells[y][x].Value {
				t.Errorf("Cells[%d][%d].Value = %d, want %d", y, x, b.Cells[y][x].Value, b2.Cells[y][x].Value)
			}
		}
	}
}
