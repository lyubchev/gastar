package main

import (
	"fmt"
	"math"
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
	x, y       int
	f, g, h    float64
	isObstacle bool
	neighbours []Cell
	previous   Cell
}

func newCell(x, y int) *Cell {
	// if it is the starting cell or the goal cell dont make it an obstacle
	isObstacle := rand.Float64() < 0.35
	if (x == 0 && y == 0) || (x == rows-1 && y == cols-1) {
		isObstacle = false
	}

	return &Cell{
		x:          x,
		y:          y,
		isObstacle: isObstacle,
		neighbours: []Cell{},
		previous:   Cell{},
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

func heuristic(cellA, cellB Cell) float64 {
	return math.Hypot(math.Abs(float64(cellA.x-cellB.x)), math.Abs(float64(cellA.y-cellB.y)))
}

func contains(elt Cell, arr []Cell) bool {
	for _, el := range arr {
		if el.x == elt.x && el.y == elt.y {
			return true
		}
	}

	return false
}

func main() {
	rl.InitWindow(windowsW+lineThickness, windowsH+lineThickness, "Gastar - A* Path Finding")

	openSet := []Cell{}
	closedSet := []Cell{}

	rl.SetTargetFPS(60)

	// Generate grid with random obstacles
	grid := [][]Cell{}
	for i := 0; i < cols; i++ {
		grid = append(grid, []Cell{})
		for j := 0; j < rows; j++ {
			grid[i] = append(grid[i], *newCell(i, j))
		}
	}

	start := grid[0][0]
	goal := grid[rows-1][cols-1]

	openSet = append(openSet, start)

	// Append all neighbours to each cell
	for i, row := range grid {
		for j := range row {
			grid[i][j].addNeighbours(grid)
		}
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		// Check if A* is still searching for path
		if len(openSet) > 0 {

			bestCell := 0
			for i, cell := range openSet {
				if cell.f < openSet[bestCell].f {
					bestCell = i
				}
			}
			current := openSet[bestCell]
			if current.x == goal.x && current.y == goal.y {
				// rl.EndDrawing()
				fmt.Println("Found the path!")
			}

			closedSet = append(closedSet, current)
			openSet = append(openSet[:bestCell], openSet[bestCell+1:]...)

			for i, neighbour := range current.neighbours {

				if contains(neighbour, closedSet) {
					continue
				}

				if !contains(neighbour, openSet) && !neighbour.isObstacle {
					current.neighbours[i].g = current.g + heuristic(neighbour, current)
					current.neighbours[i].h = heuristic(neighbour, goal)
					current.neighbours[i].f = current.neighbours[i].g + current.neighbours[i].h

				} else {

				}

			}
		}

		// Draw grid
		for _, row := range grid {
			for _, cell := range row {
				cell.draw()
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
