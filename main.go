package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	windowsW      = 500
	windowsH      = 500
	cols          = 25
	rows          = 25
	cellSize      = windowsW / cols
	LineThickness = 4
)

func main() {
	rl.InitWindow(windowsW+LineThickness, windowsH+LineThickness, "gAstar - A* path finding")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for i := 0; i < cols; i++ {
			for j := 0; j < rows; j++ {
				x := int32(i*cellSize + LineThickness)
				y := int32(j*cellSize + LineThickness)
				rl.DrawRectangle(x, y, cellSize-LineThickness, cellSize-LineThickness, rl.LightGray)
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
