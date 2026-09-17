package core

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const boardSize = 4

const winTile = 2048

type board struct {
	cells [boardSize][boardSize]int
}

type direction int

const (
	up direction = iota
	down
	left
	right
)

func slideRowLeft(row [boardSize]int) ([boardSize]int, int) {
	tiles := make([]int, 0, boardSize)
	for _, value := range row {
		if value != 0 {
			tiles = append(tiles, value)
		}
	}
	var out [boardSize]int
	gained := 0
	i, j := 0, 0
	for i < len(tiles) {
		if i+1 < len(tiles) && tiles[i] == tiles[i+1] {
			out[j] = tiles[i] * 2
			gained += tiles[i] * 2
			i += 2
		} else {
			out[j] = tiles[i]
			i++
		}
		j++
	}
	return out, gained
}

func transpose(cells [boardSize][boardSize]int) [boardSize][boardSize]int {
	var out [boardSize][boardSize]int
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			out[x][y] = cells[y][x]
		}
	}
	return out
}

func reverseRows(cells [boardSize][boardSize]int) [boardSize][boardSize]int {
	out := cells
	for y := range out {
		for x := 0; x < boardSize/2; x++ {
			out[y][x], out[y][boardSize-1-x] = out[y][boardSize-1-x], out[y][x]
		}
	}
	return out
}

func (b *board) move(dir direction) (bool, int) {
	cells := b.cells
	switch dir {
	case right:
		cells = reverseRows(cells)
	case up:
		cells = transpose(cells)
	case down:
		cells = reverseRows(transpose(cells))
	}
	gained := 0
	for y := 0; y < boardSize; y++ {
		var rowGained int
		cells[y], rowGained = slideRowLeft(cells[y])
		gained += rowGained
	}
	switch dir {
	case right:
		cells = reverseRows(cells)
	case up:
		cells = transpose(cells)
	case down:
		cells = transpose(reverseRows(cells))
	}
	moved := cells != b.cells
	b.cells = cells
	return moved, gained
}

func (b *board) emptyCells() []point {
	var empties []point
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if b.cells[y][x] == 0 {
				empties = append(empties, point{x, y})
			}
		}
	}
	return empties
}

func (b *board) spawn() {
	empties := b.emptyCells()
	if len(empties) == 0 {
		return
	}
	cell := empties[rand.Intn(len(empties))]
	value := 2
	if rand.Float64() < 0.1 {
		value = 4
	}
	b.cells[cell.y][cell.x] = value
}

func (b *board) movesAvailable() bool {
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if b.cells[y][x] == 0 {
				return true
			}
			if x+1 < boardSize && b.cells[y][x] == b.cells[y][x+1] {
				return true
			}
			if y+1 < boardSize && b.cells[y][x] == b.cells[y+1][x] {
				return true
			}
		}
	}
	return false
}

func (b *board) maxTile() int {
	max := 0
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if b.cells[y][x] > max {
				max = b.cells[y][x]
			}
		}
	}
	return max
}

var tileBackgrounds = map[int]string{
	2:    "#eee4da",
	4:    "#ede0c8",
	8:    "#f2b179",
	16:   "#f59563",
	32:   "#f67c5f",
	64:   "#f65e3b",
	128:  "#edcf72",
	256:  "#edcc61",
	512:  "#edc850",
	1024: "#edc53f",
	2048: "#edc22e",
}

func tileStyle(value int) lipgloss.Style {
	background := "#3c3a32"
	if color, ok := tileBackgrounds[value]; ok {
		background = color
	} else if value == 0 {
		background = "#cdc1b4"
	}
	foreground := "#f9f6f2"
	if value == 2 || value == 4 {
		foreground = "#776e65"
	}
	return lipgloss.NewStyle().
		Width(7).
		Align(lipgloss.Center).
		Background(lipgloss.Color(background)).
		Foreground(lipgloss.Color(foreground)).
		Bold(true)
}

func (b *board) render() string {
	var sb strings.Builder

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	width := boardSize * 7
	sb.WriteString(borderStyle.Render("╔"+strings.Repeat("═", width)+"╗") + "\n")

	for y := 0; y < boardSize; y++ {
		sb.WriteString(borderStyle.Render("║"))
		for x := 0; x < boardSize; x++ {
			value := b.cells[y][x]
			if value == 0 {
				sb.WriteString(tileStyle(0).Render(""))
			} else {
				sb.WriteString(tileStyle(value).Render(strconv.Itoa(value)))
			}
		}
		sb.WriteString(borderStyle.Render("║") + "\n")
	}

	sb.WriteString(borderStyle.Render("╚"+strings.Repeat("═", width)+"╝") + "\n")

	return sb.String()
}
