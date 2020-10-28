package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	windowsW      = 500
	windowsH      = 500
	cols          = 25
	rows          = 25
	cellSize      = windowsW / cols
	lineThickness = 4
)

type Cell struct {
	x          int
	y          int
	f          int
	g          int
	h          int
	isObstacle bool
	neighbours []Cell
}

func newCell(x, y int) *Cell {
	return &Cell{
		x:          x,
		y:          y,
		isObstacle: rand.Float64() < 0.4,
		neighbours: []Cell{},
	}
}

func (c *Cell) addNeighbours(grid [][]Cell) {
	if c.x > 0 && c.y > 0 {
		c.neighbours = append(c.neighbours, grid[c.x-1][c.y-1])
	}

	if c.y > 0 {
		c.neighbours = append(c.neighbours, grid[c.x][c.y-1])
	}

	if c.x+1 < rows && c.y > 0 {
		c.neighbours = append(c.neighbours, grid[c.x+1][c.y-1])
	}

	if c.x > 0 {
		c.neighbours = append(c.neighbours, grid[c.x-1][c.y])
	}

	if c.x+1 < rows {
		c.neighbours = append(c.neighbours, grid[c.x+1][c.y])
	}

	if c.x > 0 && c.y+1 < cols {
		c.neighbours = append(c.neighbours, grid[c.x-1][c.y+1])
	}

	if c.y+1 < cols {
		c.neighbours = append(c.neighbours, grid[c.x][c.y+1])
	}

	if c.x+1 < rows && c.y+1 < cols {
		c.neighbours = append(c.neighbours, grid[c.x+1][c.y+1])
	}
}

func (c *Cell) draw() {
	x := int32(c.x*cellSize + lineThickness)
	y := int32(c.y*cellSize + lineThickness)
	color := rl.LightGray
	if c.isObstacle {
		color = rl.Black
	}
	rl.DrawRectangle(x, y, cellSize-lineThickness, cellSize-lineThickness, color)
}

func main() {
	rl.InitWindow(windowsW+lineThickness, windowsH+lineThickness, "gAstar - A* path finding")

	rl.SetTargetFPS(60)

	grid := [][]Cell{}
	for i := 0; i < cols; i++ {
		grid = append(grid, []Cell{})
		for j := 0; j < rows; j++ {
			grid[i] = append(grid[i], *newCell(i, j))
		}
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for _, row := range grid {
			for _, cell := range row {
				cell.draw()
			}
		}
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
