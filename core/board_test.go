package core

import "testing"

func TestSlideRowLeft(t *testing.T) {
	tests := []struct {
		name   string
		row    [boardSize]int
		want   [boardSize]int
		gained int
	}{
		{"empty", [boardSize]int{0, 0, 0, 0}, [boardSize]int{0, 0, 0, 0}, 0},
		{"single tile", [boardSize]int{2, 0, 0, 0}, [boardSize]int{2, 0, 0, 0}, 0},
		{"pair merges", [boardSize]int{2, 2, 0, 0}, [boardSize]int{4, 0, 0, 0}, 4},
		{"gap pair merges", [boardSize]int{2, 0, 2, 0}, [boardSize]int{4, 0, 0, 0}, 4},
		{"double merge", [boardSize]int{2, 2, 2, 2}, [boardSize]int{4, 4, 0, 0}, 8},
		{"chain merges once each", [boardSize]int{4, 4, 8, 8}, [boardSize]int{8, 16, 0, 0}, 24},
		{"no double count", [boardSize]int{4, 4, 4, 0}, [boardSize]int{8, 4, 0, 0}, 8},
		{"full no merge", [boardSize]int{2, 4, 8, 16}, [boardSize]int{2, 4, 8, 16}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gained := slideRowLeft(tt.row)
			if got != tt.want {
				t.Errorf("slideRowLeft(%v) = %v, want %v", tt.row, got, tt.want)
			}
			if gained != tt.gained {
				t.Errorf("slideRowLeft(%v) gained = %d, want %d", tt.row, gained, tt.gained)
			}
		})
	}
}

func TestMove(t *testing.T) {
	emptyRow := [boardSize]int{}

	tests := []struct {
		name   string
		dir    direction
		start  [boardSize][boardSize]int
		want   [boardSize][boardSize]int
		moved  bool
		gained int
	}{
		{
			name:  "left merges",
			dir:   left,
			start: [boardSize][boardSize]int{{2, 2, 0, 0}, emptyRow, emptyRow, emptyRow},
			want:  [boardSize][boardSize]int{{4, 0, 0, 0}, emptyRow, emptyRow, emptyRow},
			moved: true, gained: 4,
		},
		{
			name:  "right merges",
			dir:   right,
			start: [boardSize][boardSize]int{{2, 2, 0, 0}, emptyRow, emptyRow, emptyRow},
			want:  [boardSize][boardSize]int{{0, 0, 0, 4}, emptyRow, emptyRow, emptyRow},
			moved: true, gained: 4,
		},
		{
			name:  "up merges column",
			dir:   up,
			start: [boardSize][boardSize]int{{2, 0, 0, 0}, {2, 0, 0, 0}, emptyRow, emptyRow},
			want:  [boardSize][boardSize]int{{4, 0, 0, 0}, emptyRow, emptyRow, emptyRow},
			moved: true, gained: 4,
		},
		{
			name:  "down merges column",
			dir:   down,
			start: [boardSize][boardSize]int{{2, 0, 0, 0}, {2, 0, 0, 0}, emptyRow, emptyRow},
			want:  [boardSize][boardSize]int{emptyRow, emptyRow, emptyRow, {4, 0, 0, 0}},
			moved: true, gained: 4,
		},
		{
			name:  "no change reports unmoved",
			dir:   left,
			start: [boardSize][boardSize]int{{2, 4, 8, 16}, emptyRow, emptyRow, emptyRow},
			want:  [boardSize][boardSize]int{{2, 4, 8, 16}, emptyRow, emptyRow, emptyRow},
			moved: false, gained: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &board{cells: tt.start}
			moved, gained := b.move(tt.dir)
			if moved != tt.moved {
				t.Errorf("move() moved = %v, want %v", moved, tt.moved)
			}
			if gained != tt.gained {
				t.Errorf("move() gained = %d, want %d", gained, tt.gained)
			}
			if b.cells != tt.want {
				t.Errorf("move() board = %v, want %v", b.cells, tt.want)
			}
		})
	}
}

func TestSpawn(t *testing.T) {
	b := &board{}
	b.spawn()

	count, value := 0, 0
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if b.cells[y][x] != 0 {
				count++
				value = b.cells[y][x]
			}
		}
	}
	if count != 1 {
		t.Fatalf("spawn() placed %d tiles, want exactly 1", count)
	}
	if value != 2 && value != 4 {
		t.Errorf("spawn() value = %d, want 2 or 4", value)
	}

	full := &board{}
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			full.cells[y][x] = 2
		}
	}
	before := full.cells
	full.spawn()
	if full.cells != before {
		t.Error("spawn() on a full board must change nothing")
	}
}

func TestMovesAvailable(t *testing.T) {
	tests := []struct {
		name  string
		cells [boardSize][boardSize]int
		want  bool
	}{
		{"empty board", [boardSize][boardSize]int{}, true},
		{
			name: "full board no merges",
			cells: [boardSize][boardSize]int{
				{2, 4, 8, 16},
				{32, 64, 128, 256},
				{512, 1024, 2048, 4096},
				{8192, 16384, 32768, 65536},
			},
			want: false,
		},
		{
			name: "full board with merge",
			cells: [boardSize][boardSize]int{
				{2, 4, 8, 16},
				{32, 64, 128, 256},
				{512, 1024, 2048, 4096},
				{8192, 16384, 32768, 32768},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &board{cells: tt.cells}
			if got := b.movesAvailable(); got != tt.want {
				t.Errorf("movesAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
