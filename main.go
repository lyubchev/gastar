package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	windowsW      = 500
	windowsH      = 500
	cols          = 25
	rows          = 25
	cellSize      = windowsW / cols
	lineThickness = 3
)

type Cell struct {
	x, y       int
	f, g, h    float64
	isObstacle bool
	neighbours []*Cell
	previous   *Cell
}

func newCell(x, y int) *Cell {
	// if it is the starting cell or the goal cell dont make it an obstacle
	isObstacle := rand.Float64() < 0.3
	if (x == 0 && y == 0) || (x == rows-1 && y == cols-1) {
		isObstacle = false
	}

	return &Cell{
		x:          x,
		y:          y,
		isObstacle: isObstacle,
		neighbours: []*Cell{},
		previous:   nil,
	}
}

func (c *Cell) addNeighbours(grid [][]*Cell) {
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

	if c.x+1 < rows {
		if c.y > 0 {
			c.neighbours = append(c.neighbours, grid[c.x+1][c.y-1])
		}
		c.neighbours = append(c.neighbours, grid[c.x+1][c.y])
		if c.y+1 < cols {
			c.neighbours = append(c.neighbours, grid[c.x+1][c.y+1])
		}
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

func (c *Cell) drawStep() {
	if c == nil {
		return
	}
	x := int32(c.x*cellSize + lineThickness)
	y := int32(c.y*cellSize + lineThickness)
	rl.DrawRectangle(x, y, cellSize-lineThickness, cellSize-lineThickness, rl.Lime)
}

func drawPath(lastStep *Cell) {
	path := []*Cell{}
	var temp Cell
	if lastStep != nil {
		temp = *lastStep
	}
	path = append(path, &temp)
	for temp.previous != nil {
		path = append(path, temp.previous)
		temp = *temp.previous
	}

	lastStep.drawStep()
	for _, c := range path {
		c.drawStep()
	}
}

func heuristic(cellA, cellB *Cell) float64 {
	return math.Hypot(math.Abs(float64(cellA.x-cellB.x)), math.Abs(float64(cellA.y-cellB.y)))
}

func contains(elt *Cell, arr []*Cell) bool {
	for _, el := range arr {
		if el.x == elt.x && el.y == elt.y {
			return true
		}
	}

	return false
}

func main() {
	rl.InitWindow(windowsW+lineThickness, windowsH+lineThickness, "Gastar - A* Path Finding")
	rand.Seed(time.Now().Unix())

	openSet := []*Cell{}
	closedSet := []*Cell{}

	rl.SetTargetFPS(45)

	// Generate grid with random obstacles
	grid := [][]*Cell{}
	for i := 0; i < cols; i++ {
		grid = append(grid, []*Cell{})
		for j := 0; j < rows; j++ {
			grid[i] = append(grid[i], newCell(i, j))
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

	pathFound := false

	var lastPath *Cell
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		if rl.IsKeyReleased(rl.KeyR) {
			rl.ClearBackground(rl.White)
			// Restart the loop: generate new obstacles and restart pathfinding

			// Clear the open and closed sets
			openSet = []*Cell{}
			closedSet = []*Cell{}

			// Regenerate grid with new obstacles
			for i := 0; i < cols; i++ {
				for j := 0; j < rows; j++ {
					grid[i][j] = newCell(i, j)
				}
			}

			start = grid[0][0]
			goal = grid[rows-1][cols-1]

			openSet = append(openSet, start)

			// Append all neighbours to each cell
			for i, row := range grid {
				for j := range row {
					grid[i][j].addNeighbours(grid)
				}
			}

			pathFound = false
		}

		if pathFound {
			drawPath(lastPath)
			rl.DrawText("Path found!", windowsW/2-100, windowsH/2, 35, rl.Orange)
			rl.EndDrawing()
			continue
		}

		// Draw grid
		for _, row := range grid {
			for _, cell := range row {
				cell.draw()
			}
		}

		// Check if A* is still searching for path
		if len(openSet) > 0 {

			bestCell := 0
			for i, cell := range openSet {
				if cell.f < openSet[bestCell].f {
					bestCell = i
				}
			}
			lastPath = openSet[bestCell]

			// Check if we found the path
			if lastPath.x == goal.x && lastPath.y == goal.y {
				pathFound = true
				fmt.Println("Found the path!")
				drawPath(lastPath)
				rl.EndDrawing()
				continue
			}

			closedSet = append(closedSet, lastPath)
			openSet = append(openSet[:bestCell], openSet[bestCell+1:]...)

			for i, neighbour := range lastPath.neighbours {
				if contains(neighbour, closedSet) {
					continue
				}

				tempG := lastPath.g + heuristic(neighbour, lastPath)
				newPath := false
				if !contains(neighbour, openSet) {
					if !neighbour.isObstacle {
						newPath = true
						lastPath.neighbours[i].g = tempG
						openSet = append(openSet, lastPath.neighbours[i])
					}
				} else {
					// This neighbour may come from a different current, this is why we set its new G value
					if tempG < lastPath.neighbours[i].g {
						lastPath.neighbours[i].g = tempG
						newPath = true
					}
				}

				if newPath {
					lastPath.neighbours[i].h = heuristic(neighbour, goal)
					lastPath.neighbours[i].f = lastPath.neighbours[i].g + lastPath.neighbours[i].h
					lastPath.neighbours[i].previous = lastPath
				}
			}
		} else {
			rl.DrawText("No solution!", windowsW/2, windowsH/2, 20, rl.Red)
		}

		drawPath(lastPath)
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
