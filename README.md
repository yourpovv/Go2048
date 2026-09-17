<div align="center">
  
# Go2048

**Terminal 2048 game built with Go, Bubbletea and Lipgloss.**

[![GitHub](https://img.shields.io/github/stars/yourpovv/GOOB?style=social)](https://github.com/yourpovv/GOOB)

https://github.com/user-attachments/assets/dc100a40-6fb8-4ec2-9d82-4f55974707dc

</div>

## Requirements

- Go 1.21+

## Running it

```bash
go run .
```

or build it first:

```bash
go build -o go2048.exe
.\go2048.exe
```

## Controls

- Arrow keys or WASD to slide tiles
- Q to quit
- R to restart after game over

## How it works

Slide the tiles on a 4x4 grid. Matching numbers that collide merge into their sum total giving you points. A new 2 (or 4) spawns after every move that changes the board. Reach 2048 to win
