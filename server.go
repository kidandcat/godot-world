package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/olahol/melody.v1"
)

func serverSetupRoutes() *melody.Melody {
	m := melody.New()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	m.HandleMessageBinary(func(s *melody.Session, msgb []byte) {
		fmt.Println("->", string(msgb))
		msg := strings.Split(string(msgb), ":")
		switch msg[0] {
		case "login": // -> login:user,pass
			serverLogin(s, msg[1]) // <- player:{player}
		case "create_mesh":
			serverCreateMesh(m, s, msg[1])
		case "world_around":
			serverSendWorldAround(s, msg[1])
		case "walk_to":
			serverWalkTo(s, msg[1])
		case "notify_movement":
			serverNotifyMovement(s, msg[1])
		default:
			fmt.Println("Unknown command", msg)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		fmt.Println("New connection")
	})
	return m
}

/* -------- */
/* COMMANDS */
/* -------- */
/*
Available client commands:
- login:ok
- player:<x>,<z>
- newmesh:<type>,<x>,<z>,<rotation>
*/

func serverLogin(s *melody.Session, data string) {
	login := strings.Split(data, ",")
	p := &tPlayer{}
	if login[0] == "jairo" && login[1] == "new" { // TODO proper errors
		p = &tPlayer{
			Name:     "Jairo",
			Position: tPos{X: 0, Z: 0},
		}
		db.newPlayer(p)
		// send player to client
		send(s, "login:ok")
		// save player in socket
	} else if login[1] == "get" {
		// TODO receive player name
		p = db.getPlayer("1")
		x := strconv.FormatInt(p.Position.X, 10)
		z := strconv.FormatInt(p.Position.Z, 10)
		// send player to client
		send(s, "player:"+x+","+z)
	}
	// save player in socket
	s.Set("player", p.ID)
}

func serverCreateMesh(m *melody.Melody, s *melody.Session, d string) {
	data := strings.Split(d, ",")
	t, _ := strconv.ParseInt(data[0], 10, 64)
	x, _ := strconv.ParseInt(data[1], 10, 64)
	z, _ := strconv.ParseInt(data[2], 10, 64)
	y, _ := strconv.ParseInt(data[3], 10, 64)
	mesh := &tMesh{
		Type: t,
		Position: tPos{
			X: x,
			Z: z,
		},
		VerticalLevel: y,
		Walkable:      true,
		WalkingCost:   1,
		Rotation:      data[4],
	}
	db.newMesh(mesh)
	broadcast(m, s, "newmesh:"+data[0]+","+data[1]+","+data[2]+","+data[3]+","+data[4])
}

func serverSendWorldAround(s *melody.Session, d string) {
	data := strings.Split(d, ",")
	x, _ := strconv.ParseInt(data[0], 10, 64)
	z, _ := strconv.ParseInt(data[1], 10, 64)
	meshes := []string{}
	db.getNearbyMeshes(x, z, 10, &meshes, nil) // TODO calculate proper screen distance
	m := &tMesh{}
	for _, v := range meshes {
		db.getMesh(m, v)
		t := strconv.FormatInt(m.Type, 10)
		x := strconv.FormatInt(m.Position.X, 10)
		z := strconv.FormatInt(m.Position.Z, 10)
		y := strconv.FormatInt(m.VerticalLevel, 10)
		send(s, "newmesh:"+t+","+x+","+z+","+y+","+m.Rotation)
	}
}

func serverWalkTo(s *melody.Session, d string) {
	data := strings.Split(d, ",")
	x, _ := strconv.ParseInt(data[0], 10, 64)
	z, _ := strconv.ParseInt(data[1], 10, 64)
	pIDi, _ := s.Get("player")
	pID := pIDi.(string)

	player := db.getPlayer(pID)

	origin := db.getMeshByPos(player.Position.X, player.Position.Z, 0)
	destination := db.getMeshByPos(x, z, 0)

	// do nothing if origin or destination don't exist or it's the same position
	if origin == nil || destination == nil || (origin.Position.X == destination.Position.X && origin.Position.Z == destination.Position.Z) {
		return
	}
	steps := movementCalculatePath(destination, origin)
	steps = steps[1:] // remove first step (it is player's current position)
	for _, step := range steps {
		send(s, "move:"+strconv.FormatInt(step.Position.X, 10)+","+strconv.FormatInt(step.Position.Z, 10))
	}
}

func serverNotifyMovement(s *melody.Session, d string) {
	// TODO check movement validity
	// TODO update player's world on movement
	data := strings.Split(d, ",")
	x, _ := strconv.ParseInt(data[0], 10, 64)
	z, _ := strconv.ParseInt(data[1], 10, 64)
	pIDi, _ := s.Get("player")
	pID := pIDi.(string)

	player := db.getPlayer(pID)

	player.Position.X = x
	player.Position.Z = z
	db.setPlayerPos(player)
}

func broadcast(m *melody.Melody, s *melody.Session, d string) {
	fmt.Println("<--", d)
	m.Broadcast([]byte(d))
}

func send(s *melody.Session, d string) {
	fmt.Println("<-", d)
	s.Write([]byte(d))
}
