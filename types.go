package main

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
	VerticalLevel int64
	Walkable      bool
	WalkingCost   int
}
