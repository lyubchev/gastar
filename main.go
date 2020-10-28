package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	windowsW      = 500
	windowsH      = 500
	cols          = 25
	rows          = 25
	cellSize      = windowsW / cols
	lineThickness = 4
)

func main() {
	rl.InitWindow(windowsW+lineThickness, windowsH+lineThickness, "gAstar - A* path finding")

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for i := 0; i < cols; i++ {
			for j := 0; j < rows; j++ {
				x := int32(i*cellSize + lineThickness)
				y := int32(j*cellSize + lineThickness)
				rl.DrawRectangle(x, y, cellSize-lineThickness, cellSize-lineThickness, rl.LightGray)
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
