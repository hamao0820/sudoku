package solver

func Backtrack(b *Board) bool {
	if b.IsSolved() {
		return true
	}
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if !b.IsEmpty(y, x) {
				continue
			}

			for value := 9; value > 0; value-- {
				if !b.IsLegal(y, x, value) {
					continue
				}

				b.Cells[y][x].Value = value
				if Backtrack(b) {
					return true
				}
				b.Cells[y][x].Value = 0
			}
			return false
		}
	}
	return false
}
