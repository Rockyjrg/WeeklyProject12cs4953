package main

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rectangle struct {
	X, Y, Width, Height int32
}

type Room struct {
	X, Y, Width, Height int32
}

type Node struct {
	Rect        Rectangle
	Left, Right *Node
	Room        *Room
	Color       rl.Color
}

func Split(node *Node, minSize int32) {
	if node == nil {
		return
	}

	node.Color = rl.Color{
		R: uint8(rand.IntN(255)),
		G: uint8(rand.IntN(255)),
		B: uint8(rand.IntN(255)),
		A: 255,
	}

	rect := node.Rect

	if rect.Width < 2*minSize && rect.Height < 2*minSize {
		return
	}

	splitHorizontally := rand.IntN(2) == 0

	if rect.Width > rect.Height && rect.Width/rect.Height >= 1 {
		splitHorizontally = false
	} else if rect.Height > rect.Width && rect.Height/rect.Width >= 1 {
		splitHorizontally = true
	}

	if splitHorizontally {
		max := rect.Height - minSize
		if max <= minSize {
			return
		}

		split := rand.Int32N(max-minSize) + minSize

		node.Left = &Node{Rect: Rectangle{X: rect.X, Y: rect.Y, Width: rect.Width, Height: split}}
		node.Right = &Node{Rect: Rectangle{X: rect.X, Y: rect.Y + split, Width: rect.Width, Height: rect.Height - split}}
	} else {
		max := rect.Width - minSize
		if max <= minSize {
			return
		}

		split := rand.Int32N(max-minSize) + minSize

		node.Left = &Node{Rect: Rectangle{X: rect.X, Y: rect.Y, Width: split, Height: rect.Height}}
		node.Right = &Node{Rect: Rectangle{X: rect.X + split, Y: rect.Y, Width: rect.Width - split, Height: rect.Height}}
	}

	Split(node.Left, minSize)
	Split(node.Right, minSize)
}

func CreateRooms(node *Node) {
	if node.Left != nil || node.Right != nil {

		// not a leaf node
		if node.Left != nil {
			CreateRooms(node.Left)
		}
		if node.Right != nil {
			CreateRooms(node.Right)
		}

		return
	}

	// leaf node
	padding := int32(4)

	roomWidth := rand.Int32N(node.Rect.Width/2) + padding
	roomHeight := rand.Int32N(node.Rect.Height/2) + padding

	roomX := node.Rect.X + rand.Int32N(node.Rect.Width-roomWidth)
	roomY := node.Rect.Y + rand.Int32N(node.Rect.Height-roomHeight)

	node.Room = &Room{X: roomX, Y: roomY, Width: roomWidth, Height: roomHeight}
}

func ConnectRooms(node *Node) {
	if node.Left != nil && node.Right != nil {
		centerLeft := GetRoomCenter(node.Left)
		centerRight := GetRoomCenter(node.Right)

		// Hallway: first horizontal, then vertical
		if rand.IntN(2) == 0 {
			// Horizontal
			corridors = append(corridors, Rectangle{X: min(centerLeft.X, centerRight.X), Y: centerLeft.Y, Width: abs(centerLeft.X - centerRight.X), Height: 4})
			corridors = append(corridors, Rectangle{X: centerRight.X, Y: min(centerLeft.Y, centerRight.Y), Width: 4, Height: abs(centerLeft.Y - centerRight.Y)})
		} else {
			// Vertical
			corridors = append(corridors, Rectangle{X: centerLeft.X, Y: min(centerLeft.Y, centerRight.Y), Width: 4, Height: abs(centerLeft.Y - centerRight.Y)})
			corridors = append(corridors, Rectangle{X: min(centerLeft.X, centerRight.X), Y: centerRight.Y, Width: abs(centerLeft.X - centerRight.X), Height: 4})
		}

		ConnectRooms(node.Left)
		ConnectRooms(node.Right)
	}
}

func GetRoomCenter(node *Node) (center Rectangle) {
	if node.Room != nil {
		return Rectangle{X: node.Room.X + node.Room.Width/2, Y: node.Room.Y + node.Room.Height/2}
	} else if node.Left != nil {
		return GetRoomCenter(node.Left)
	} else if node.Right != nil {
		return GetRoomCenter(node.Right)
	}
	return Rectangle{X: 0, Y: 0}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func abs(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func DrawDungeon(root *Node) {
	if root == nil {
		return
	}

	if root.Left == nil && root.Right == nil {
		// Draw the room
		if root.Room != nil {
			rl.DrawRectangle(root.Room.X, root.Room.Y, root.Room.Width, root.Room.Height, root.Color)
		}
	} else {
		// Draw partitions
		rl.DrawRectangleLines(root.Rect.X, root.Rect.Y, root.Rect.Width, root.Rect.Height, root.Color)

		DrawDungeon(root.Left)
		DrawDungeon(root.Right)
	}
}

// Draw hallways separately
func DrawCorridors() {
	for _, c := range corridors {
		rl.DrawRectangle(c.X, c.Y, c.Width, c.Height, rl.Brown)
	}
}
