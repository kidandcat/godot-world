package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"worldserver/structs"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

var db *DB

func main() {
	db = initializeDB()
	r := gin.Default()
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessageBinary(func(s *melody.Session, msgb []byte) {
		msg := strings.Split(string(msgb), ":")
		switch msg[0] {
		case "login": // -> login:user,pass
			handlerLogin(s, msg[1]) // <- player:{player}
		case "create_mesh":
			handlerCreateMesh(s, msg[1])
		case "world_around":
			sendWorldAround(s, msg[1])
		default:
			fmt.Println("Unknown command", msg)
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		fmt.Println("- - - new connection")
	})

	r.Run(":5000")
}

func handlerLogin(s *melody.Session, data string) {
	login := strings.Split(data, ",")
	p := &structs.TPlayer{}
	if login[0] == "jairo" && login[1] == "new" {
		p = &structs.TPlayer{
			Name:     "Jairo",
			Position: structs.TPos{X: 0, Z: 0},
		}
		db.newPlayer(p)
		// send player to client
		s.Write([]byte("login:ok"))
		// save player in socket
	} else if login[1] == "get" {
		db.getPlayer(p, "1")
		jp, e := p.MarshalJSON()
		check(e)
		// send player to client
		s.Write([]byte("player:" + string(jp)))
		// save player in socket
	}
	fmt.Printf("player %+v\n", p)
	s.Set("player", p.ID)
}

func handlerCreateMesh(s *melody.Session, d string) {
	data := strings.Split(d, ",")
	t, _ := strconv.ParseInt(data[0], 10, 64)
	x, _ := strconv.ParseInt(data[1], 10, 64)
	z, _ := strconv.ParseInt(data[2], 10, 64)
	mesh := &structs.TMesh{
		Type: t,
		Position: structs.TPos{
			X: x,
			Z: z,
		},
	}
	db.newMesh(mesh)
	fmt.Println("creating mesh", mesh)
}

func sendWorldAround(s *melody.Session, d string) {
	data := strings.Split(d, ",")
	x, _ := strconv.ParseInt(data[0], 10, 64)
	z, _ := strconv.ParseInt(data[1], 10, 64)
	meshes := []string{}
	db.getNearbyMeshes(x, z, 10, &meshes)
	fmt.Println("found meshes for", x, z, meshes)
	m := &structs.TMesh{}
	for _, v := range meshes {
		db.getMesh(m, v)
		t := strconv.FormatInt(m.Type, 10)
		x := strconv.FormatInt(m.Position.X, 10)
		z := strconv.FormatInt(m.Position.Z, 10)
		s.Write([]byte("newmesh:" + t + "," + x + "," + z))
	}
}
