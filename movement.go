package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/beefsack/go-astar"
)

var cache map[string]*tMesh

func cachedGetMeshByPos(x, z, y int64) *tMesh {
	sx := strconv.FormatInt(x, 10)
	sz := strconv.FormatInt(z, 10)
	sy := strconv.FormatInt(y, 10)

	tMesh := cache[sx+"-"+sz+"-"+sy]
	if tMesh == nil { // if its not found, cache it
		tMesh = db.getMeshByPos(x, z, y)
		cache[sx+"-"+sz+"-"+sy] = tMesh
	}
	return tMesh
}

func movementCalculatePath(a, b *tMesh) []*tMesh {
	Ax := strconv.FormatInt(a.Position.X, 10)
	Az := strconv.FormatInt(a.Position.Z, 10)
	Ay := strconv.FormatInt(a.VerticalLevel, 10)
	Bx := strconv.FormatInt(b.Position.X, 10)
	Bz := strconv.FormatInt(b.Position.Z, 10)
	By := strconv.FormatInt(b.VerticalLevel, 10)
	cache = map[string]*tMesh{
		Ax + "-" + Az + "-" + Ay: a,
		Bx + "-" + Bz + "-" + By: b,
	}
	res := []*tMesh{}
	fmt.Println("--> movementCalculatePath", a.Position, "->", b.Position)
	path, distance, found := astar.Path(a, b)
	fmt.Println("<-- astar.Path:", path, "distance:", distance, "found:", found)
	for _, m := range path {
		res = append(res, m.(*tMesh))
	}
	return res
}

func (t *tMesh) PathNeighbors() []astar.Pather {
	res := []astar.Pather{}
	log := []tPos{}
	for _, m := range []*tMesh{
		cachedGetMeshByPos(t.Position.X, t.Position.Z-1, t.VerticalLevel), // UP is z-1
		cachedGetMeshByPos(t.Position.X+1, t.Position.Z, t.VerticalLevel), // RIGHT is x+1
		cachedGetMeshByPos(t.Position.X, t.Position.Z+1, t.VerticalLevel), // DOWN is z+1
		cachedGetMeshByPos(t.Position.X-1, t.Position.Z, t.VerticalLevel), // LEFT is x-1
	} {
		if m != nil {
			res = append(res, m)
			log = append(log, m.Position)
		}
	}
	return res
}

func (t *tMesh) PathNeighborCost(to astar.Pather) float64 {
	toMesh := to.(*tMesh)
	return float64(toMesh.WalkingCost)
}

func (t *tMesh) PathEstimatedCost(to astar.Pather) float64 {
	toMesh := to.(*tMesh)
	distance := math.Abs(float64(toMesh.Position.X)-float64(t.Position.X)) + math.Abs(float64(toMesh.Position.Z)-float64(t.Position.Z))
	return distance
}
