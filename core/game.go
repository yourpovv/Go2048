package core

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameState int

const (
	playing gameState = iota
	gameOver
)

type Game struct {
	board *board
	score int
	won   bool
	state gameState
}

func NewGame() *Game {
	g := &Game{board: &board{}, state: playing}
	g.board.spawn()
	g.board.spawn()
	return g
}

func (g *Game) Init() tea.Cmd {
	return nil
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "r":
			if g.state == gameOver {
				return NewGame(), nil
			}
		case "left", "a", "h":
			g.move(left)
		case "right", "d", "l":
			g.move(right)
		case "up", "w", "k":
			g.move(up)
		case "down", "s", "j":
			g.move(down)
		}
	}

	return g, nil
}

func (g *Game) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true).
		MarginBottom(1).
		Render("🔢 Go2048")

	scoreText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Render(fmt.Sprintf("Score: %d", g.score))

	gameBoard := g.board.render()

	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		MarginTop(1).
		Render("Arrows/WASD: Move | Q: Quit | R: Restart")

	var statusBar string
	switch g.state {
	case gameOver:
		statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render("💀 GAME OVER | Press R to restart")
	default:
		if g.won {
			statusBar = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFD700")).
				Bold(true).
				Render("🎉 2048! Keep going")
		} else {
			statusBar = ""
		}
	}

	view := fmt.Sprintf("%s\n%s\n\n%s\n%s", title, scoreText, gameBoard, controls)

	if statusBar != "" {
		view = fmt.Sprintf("%s\n\n%s", view, statusBar)
	}

	return lipgloss.NewStyle().
		MarginLeft(2).
		MarginTop(1).
		Render(view)
}

func (g *Game) move(dir direction) {
	if g.state != playing {
		return
	}
	moved, gained := g.board.move(dir)
	if !moved {
		return
	}
	g.score += gained
	g.board.spawn()
	if g.board.maxTile() >= winTile {
		g.won = true
	}
	if !g.board.movesAvailable() {
		g.state = gameOver
	}
}
