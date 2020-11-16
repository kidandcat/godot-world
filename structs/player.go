package structs

const (
	NodeTypePlayer = "player"
	NodeTypeMesh   = "mesh"
)

//easyjson:json
type TPlayer struct {
	ID       string
	Name     string
	Position TPos
}

type TPos struct {
	X int64
	Z int64
}

type TMesh struct {
	ID       string
	Type     int64
	Position TPos
}
