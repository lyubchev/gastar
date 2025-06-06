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
	controlPanelH = 60
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

func newCell(x, y int, obstacleDensity float64) *Cell {
	// if it is the starting cell or the goal cell dont make it an obstacle
	isObstacle := rand.Float64() < obstacleDensity
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
	y := int32(c.y*cellSize + lineThickness + controlPanelH)
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
	y := int32(c.y*cellSize + lineThickness + controlPanelH)
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

func drawButton(x, y, width, height int32, text string, pressed bool) bool {
	color := rl.LightGray
	if pressed {
		color = rl.Gray
	}
	
	rl.DrawRectangle(x, y, width, height, color)
	rl.DrawRectangleLines(x, y, width, height, rl.Black)
	
	textWidth := rl.MeasureText(text, 16)
	textX := x + (width-textWidth)/2
	textY := y + (height-16)/2
	rl.DrawText(text, textX, textY, 16, rl.Black)
	
	mousePos := rl.GetMousePosition()
	return rl.IsMouseButtonPressed(rl.MouseLeftButton) &&
		mousePos.X >= float32(x) && mousePos.X <= float32(x+width) &&
		mousePos.Y >= float32(y) && mousePos.Y <= float32(y+height)
}

func drawSlider(x, y, width int32, value, min, max float64, label string) float64 {
	height := int32(20)
	
	// Draw slider track
	rl.DrawRectangle(x, y, width, height, rl.LightGray)
	rl.DrawRectangleLines(x, y, width, height, rl.Black)
	
	// Draw slider handle
	handleX := x + int32(float64(width-10)*((value-min)/(max-min)))
	rl.DrawRectangle(handleX, y-5, 10, height+10, rl.DarkGray)
	
	// Draw label and value
	rl.DrawText(label, x, y-20, 14, rl.Black)
	valueText := fmt.Sprintf("%.2f", value)
	rl.DrawText(valueText, x+width-rl.MeasureText(valueText, 14), y-20, 14, rl.Black)
	
	// Handle mouse interaction
	mousePos := rl.GetMousePosition()
	if rl.IsMouseButtonDown(rl.MouseLeftButton) &&
		mousePos.X >= float32(x) && mousePos.X <= float32(x+width) &&
		mousePos.Y >= float32(y-5) && mousePos.Y <= float32(y+height+5) {
		newValue := min + (max-min)*float64(mousePos.X-float32(x))/float64(width)
		if newValue < min {
			newValue = min
		}
		if newValue > max {
			newValue = max
		}
		return newValue
	}
	
	return value
}

func generateGrid(obstacleDensity float64) [][]*Cell {
	grid := [][]*Cell{}
	for i := 0; i < cols; i++ {
		grid = append(grid, []*Cell{})
		for j := 0; j < rows; j++ {
			grid[i] = append(grid[i], newCell(i, j, obstacleDensity))
		}
	}
	
	// Add neighbours to each cell
	for i, row := range grid {
		for j := range row {
			grid[i][j].addNeighbours(grid)
		}
	}
	
	return grid
}

func resetPathfinding(grid [][]*Cell) ([]*Cell, []*Cell, *Cell, *Cell) {
	openSet := []*Cell{}
	closedSet := []*Cell{}
	start := grid[0][0]
	goal := grid[rows-1][cols-1]
	
	// Reset all cell values for pathfinding
	for i := range grid {
		for j := range grid[i] {
			grid[i][j].f = 0
			grid[i][j].g = 0
			grid[i][j].h = 0
			grid[i][j].previous = nil
		}
	}
	
	openSet = append(openSet, start)
	return openSet, closedSet, start, goal
}

func main() {
	rl.InitWindow(windowsW+lineThickness, windowsH+lineThickness+controlPanelH, "Gastar - A* Path Finding")
	rand.Seed(time.Now().Unix())

	// Control variables
	obstacleDensity := 0.3
	speed := 1.0
	targetFPS := int32(30 + speed*30) // 30-60 FPS based on speed
	
	rl.SetTargetFPS(targetFPS)

	// Generate initial grid
	grid := generateGrid(obstacleDensity)
	openSet, closedSet, _, goal := resetPathfinding(grid)

	pathFound := false
	var lastPath *Cell
	stepCounter := 0
	
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		// Draw control panel background
		rl.DrawRectangle(0, 0, windowsW+lineThickness, controlPanelH, rl.RayWhite)
		rl.DrawLine(0, controlPanelH, windowsW+lineThickness, controlPanelH, rl.Black)

		// Draw restart button
		if drawButton(10, 10, 80, 40, "Restart", false) {
			grid = generateGrid(obstacleDensity)
			openSet, closedSet, _, goal = resetPathfinding(grid)
			pathFound = false
			lastPath = nil
			stepCounter = 0
		}

		// Draw density slider
		newDensity := drawSlider(110, 25, 120, obstacleDensity, 0.0, 0.8, "Density")
		if newDensity != obstacleDensity {
			obstacleDensity = newDensity
			grid = generateGrid(obstacleDensity)
			openSet, closedSet, _, goal = resetPathfinding(grid)
			pathFound = false
			lastPath = nil
			stepCounter = 0
		}

		// Draw speed slider
		newSpeed := drawSlider(250, 25, 120, speed, 0.1, 2.0, "Speed")
		if newSpeed != speed {
			speed = newSpeed
			targetFPS = int32(15 + speed*45) // 15-60 FPS based on speed
			rl.SetTargetFPS(targetFPS)
		}

		// Draw status text
		statusText := "Searching..."
		if pathFound {
			statusText = "Path Found!"
		} else if len(openSet) == 0 {
			statusText = "No Solution!"
		}
		rl.DrawText(statusText, 390, 25, 16, rl.DarkGreen)

		// Handle keyboard restart (keep the 'R' key functionality)
		if rl.IsKeyReleased(rl.KeyR) {
			grid = generateGrid(obstacleDensity)
			openSet, closedSet, _, goal = resetPathfinding(grid)
			pathFound = false
			lastPath = nil
			stepCounter = 0
		}

		// Draw grid
		for _, row := range grid {
			for _, cell := range row {
				cell.draw()
			}
		}

		// Draw start and goal markers
		rl.DrawRectangle(lineThickness, controlPanelH+lineThickness, cellSize-lineThickness, cellSize-lineThickness, rl.Green)
		rl.DrawRectangle(int32((rows-1)*cellSize+lineThickness), int32((cols-1)*cellSize+lineThickness+controlPanelH), cellSize-lineThickness, cellSize-lineThickness, rl.Red)

		if pathFound {
			drawPath(lastPath)
			rl.DrawText("Path found!", windowsW/2-100, windowsH/2+controlPanelH, 35, rl.Orange)
			rl.EndDrawing()
			continue
		}

		// A* algorithm execution with speed control
		stepsPerFrame := int(speed * 5) // More steps per frame for higher speed
		if stepsPerFrame < 1 {
			stepsPerFrame = 1
		}

		for step := 0; step < stepsPerFrame && len(openSet) > 0 && !pathFound; step++ {
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
				break
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
			stepCounter++
		}

		if len(openSet) == 0 && !pathFound {
			rl.DrawText("No solution!", windowsW/2-60, windowsH/2+controlPanelH, 20, rl.Red)
		}

		// Draw current path being explored
		if lastPath != nil {
			drawPath(lastPath)
		}

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
