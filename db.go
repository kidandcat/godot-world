package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"worldserver/structs"

	"github.com/tidwall/buntdb"
)

type DB struct {
	db *buntdb.DB
}

func initializeDB() *DB {
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}

	db.CreateIndex("names", "player:*", buntdb.IndexString)
	db.CreateSpatialIndex("playerPos", "player:*:pos", buntdb.IndexRect)
	db.CreateSpatialIndex("meshPos", "mesh:*:pos", buntdb.IndexRect)

	d := &DB{
		db: db,
	}

	d.ensureIndexExists()

	return d
}

/* ------ */
/* PLAYER */
/* ------ */

func (db *DB) newPlayer(p *structs.TPlayer) {
	if p == nil {
		panic("newPlayer: Player is NIL")
	}
	err := db.db.Update(func(tx *buntdb.Tx) error {
		index := nextIndex(tx)
		tx.Set("player:"+index+":name", p.Name, nil)
		tx.Set("player:"+index+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		p.ID = index
		return nil
	})
	check(err)
}

func (db *DB) setPlayerName(p *structs.TPlayer) {
	err := db.db.Update(func(tx *buntdb.Tx) error {
		tx.Set("player:"+p.ID+":name", p.Name, nil)
		return nil
	})
	check(err)
}

func (db *DB) setPlayerPos(p *structs.TPlayer) {
	err := db.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("player:"+p.ID+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		return err
	})
	check(err)
}

func (db *DB) getPlayerName(p *structs.TPlayer) {
	db.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("player:" + p.ID + ":name")
		if err != nil {
			return err
		}
		p.Name = val
		fmt.Printf("name is %s\n", val)
		return nil
	})
}

func (db *DB) getPlayerPos(p *structs.TPlayer) {
	db.db.View(func(tx *buntdb.Tx) error {
		pos, err := tx.Get("player:" + p.ID + ":pos")
		if err != nil {
			return err
		}
		var X, Z int64
		fmt.Sscanf(pos, "[%d %d]", &X, &Z)
		p.Position.X = X
		p.Position.Z = Z
		return nil
	})
}

func (db *DB) getPlayer(p *structs.TPlayer, ID string) {
	db.db.View(func(tx *buntdb.Tx) error {
		name, err := tx.Get("player:" + ID + ":name")
		if err != nil {
			return err
		}
		pos, err := tx.Get("player:" + ID + ":pos")
		if err != nil {
			return err
		}
		var X, Z int64
		fmt.Sscanf(pos, "[%d %d]", &X, &Z)
		p.ID = ID
		p.Name = name
		p.Position.X = X
		p.Position.Z = Z
		return nil
	})
}

/* ---- */
/* MESH */
/* ---- */

func (db *DB) newMesh(p *structs.TMesh) {
	if p == nil {
		panic("newMesh: TMesh is NIL")
	}
	err := db.db.Update(func(tx *buntdb.Tx) error {
		oldMeshes := []string{}
		X := strconv.FormatInt(p.Position.X, 10)
		Z := strconv.FormatInt(p.Position.Z, 10)
		tx.Intersects("meshPos", "["+X+" "+Z+"],["+X+" "+Z+"]", func(key, val string) bool {
			fmt.Println("mesh to delete", key, val)
			oldMeshes = append(oldMeshes, key)
			return true
		})
		for _, m := range oldMeshes {
			_, err := tx.Delete(m)
			if err != nil {
				fmt.Println("error deleting", m, err)
			}
		}

		index := nextIndex(tx)
		tx.Set("mesh:"+index+":type", strconv.FormatInt(p.Type, 10), nil)
		tx.Set("mesh:"+index+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		p.ID = index
		return nil
	})
	check(err)
}

func (db *DB) setMesh(p *structs.TMesh) {
	err := db.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("mesh:"+p.ID+":type", strconv.FormatInt(p.Type, 10), nil)
		if err != nil {
			return err
		}
		_, _, err = tx.Set("mesh:"+p.ID+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		return err
	})
	check(err)
}

func (db *DB) getMesh(p *structs.TMesh, ID string) {
	db.db.View(func(tx *buntdb.Tx) error {
		T, err := tx.Get("mesh:" + ID + ":type")
		check(err)
		pos, err := tx.Get("mesh:" + ID + ":pos")
		check(err)
		var X, Z int64
		fmt.Sscanf(pos, "[%d %d]", &X, &Z)
		p.ID = ID
		p.Type, err = strconv.ParseInt(T, 10, 64)
		p.Position.X = X
		p.Position.Z = Z
		check(err)
		return nil
	})
}

func (db *DB) getNearbyMeshes(x, y, dist int64, meshes *[]string) {
	topLeftX := strconv.FormatInt(x-dist, 10)
	topLeftY := strconv.FormatInt(y-dist, 10)
	bottomRightX := strconv.FormatInt(x+dist, 10)
	bottomRightY := strconv.FormatInt(y+dist, 10)
	fmt.Println("[" + topLeftX + " " + topLeftY + "],[" + bottomRightX + " " + bottomRightY + "]")
	db.db.View(func(tx *buntdb.Tx) error {
		tx.Intersects("meshPos", "["+topLeftX+" "+topLeftY+"],["+bottomRightX+" "+bottomRightY+"]", func(key, val string) bool {
			kID := strings.Split(key, ":")[1]
			*meshes = append(*meshes, kID)
			return true
		})
		return nil
	})
}

/* ----- */
/* UTILS */
/* ----- */

func (db *DB) ensureIndexExists() {
	err := db.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Get("index")
		if err != nil {
			_, _, err := tx.Set("index", "0", nil)
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func nextIndex(tx *buntdb.Tx) string {
	index, e := tx.Get("index")
	check(e)
	iindex, e := strconv.ParseInt(index, 10, 64)
	check(e)
	next := iindex + 1
	nextIndexStr := fmt.Sprintf("%v", next)
	tx.Set("index", nextIndexStr, nil)
	return nextIndexStr
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
