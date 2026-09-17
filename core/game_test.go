package core

import "testing"

func TestWin(t *testing.T) {
	g := NewGame()
	g.board.cells = [boardSize][boardSize]int{{1024, 1024, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}
	g.move(left)

	if !g.won {
		t.Error("expected won after reaching 2048")
	}
	if g.score != 2048 {
		t.Errorf("score = %d, want 2048", g.score)
	}
}

func TestGameOver(t *testing.T) {
	g := NewGame()
	g.board.cells = [boardSize][boardSize]int{
		{8, 16, 32, 64},
		{128, 256, 512, 1024},
		{2048, 4096, 8192, 16384},
		{0, 32768, 65536, 131072},
	}
	g.move(left)

	if g.state != gameOver {
		t.Errorf("state = %v, want gameOver after the last move", g.state)
	}
}

func TestNewGame(t *testing.T) {
	g := NewGame()

	count := 0
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if g.board.cells[y][x] != 0 {
				count++
			}
		}
	}
	if count != 2 {
		t.Errorf("NewGame spawned %d tiles, want 2", count)
	}
	if g.state != playing {
		t.Errorf("state = %v, want playing", g.state)
	}
}
