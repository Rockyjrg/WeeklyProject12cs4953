package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var root *Node
var corridors []Rectangle

func generateDungeon() {
	root = &Node{Rect: Rectangle{X: 0, Y: 0, Width: 1920, Height: 1080}}
	corridors = []Rectangle{}
	Split(root, 80)
	CreateRooms(root)
	ConnectRooms(root)
}

func main() {
	screenWidth := int32(1920)
	screenHeight := int32(1080)

	rl.InitWindow(screenWidth, screenHeight, "BSP Dungeon Generator")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	generateDungeon()

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyR) {
			generateDungeon()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		DrawDungeon(root)
		DrawCorridors()

		rl.EndDrawing()
	}
}
