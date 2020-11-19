package main

const (
	rotationUp    = "up"
	rotationRight = "right"
	rotationDown  = "down"
	rotationLeft  = "left"
)

type tPlayer struct {
	ID       string
	Name     string
	Position tPos
}

type tPos struct {
	X int64
	Z int64
}

type tMesh struct {
	ID            string
	Type          int64
	Position      tPos
	Rotation      string
	VerticalLevel int64
	Walkable      bool
	WalkingCost   int
}
