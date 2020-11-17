package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/buntdb"
)

/* ---- */
/* MESH */
/* ---- */

func (db *DB) newMesh(p *tMesh) {
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

		index := dbNextIndex(tx)
		tx.Set("mesh:"+index+":type", strconv.FormatInt(p.Type, 10), nil)
		tx.Set("mesh:"+index+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		tx.Set("mesh:"+index+":verticalLevel", strconv.FormatInt(p.VerticalLevel, 10), nil)
		tx.Set("mesh:"+index+":walkable", utilsBoolToString(p.Walkable), nil)
		tx.Set("mesh:"+index+":walkingCost", strconv.Itoa(p.WalkingCost), nil)
		p.ID = index
		return nil
	})
	utilsCheck(err)
}

func (db *DB) getMesh(p *tMesh, ID string) {
	db.db.View(func(tx *buntdb.Tx) error {
		T, err := tx.Get("mesh:" + ID + ":type")
		utilsCheck(err)
		pos, err := tx.Get("mesh:" + ID + ":pos")
		utilsCheck(err)
		verticalLevel, err := tx.Get("mesh:" + ID + ":verticalLevel")
		utilsCheck(err)
		walkable, err := tx.Get("mesh:" + ID + ":walkable")
		utilsCheck(err)
		walkingCost, err := tx.Get("mesh:" + ID + ":walkingCost")
		utilsCheck(err)
		var X, Z int64
		fmt.Sscanf(pos, "[%d %d]", &X, &Z)
		p.ID = ID
		p.Type, err = strconv.ParseInt(T, 10, 64)
		p.Position.X = X
		p.Position.Z = Z
		p.VerticalLevel, err = strconv.ParseInt(verticalLevel, 10, 64)
		p.Walkable = utilsStringToBool(walkable)
		p.WalkingCost, err = strconv.Atoi(walkingCost)
		utilsCheck(err)
		return nil
	})
}

func (db *DB) getNearbyMeshes(x, z, dist int64, meshesKeys *[]string, meshesValues *[]string) {
	if meshesValues == nil {
		meshesValues = &[]string{}
	}
	topLeftX := strconv.FormatInt(x-dist, 10)
	topLeftZ := strconv.FormatInt(z-dist, 10)
	bottomRightX := strconv.FormatInt(x+dist, 10)
	bottomRightZ := strconv.FormatInt(z+dist, 10)
	fmt.Println("NEARBY [" + topLeftX + " " + topLeftZ + "],[" + bottomRightX + " " + bottomRightZ + "]")
	db.db.View(func(tx *buntdb.Tx) error {
		tx.Intersects("meshPos", "["+topLeftX+" "+topLeftZ+"],["+bottomRightX+" "+bottomRightZ+"]", func(key, val string) bool {
			kID := strings.Split(key, ":")[1]
			*meshesKeys = append(*meshesKeys, kID)
			*meshesValues = append(*meshesValues, val)
			return true
		})
		return nil
	})
}

func (db *DB) getMeshByPos(x, z, y int64) *tMesh {
	res := &tMesh{}
	meshesKeys := []string{}
	_x := strconv.FormatInt(x, 10)
	_z := strconv.FormatInt(z, 10)
	db.db.View(func(tx *buntdb.Tx) error {
		return tx.Intersects("meshPos", "["+_x+" "+_z+"],["+_x+" "+_z+"]", func(key, val string) bool {
			kID := strings.Split(key, ":")[1]
			meshesKeys = append(meshesKeys, kID)
			return true
		})
	})
	for _, m := range meshesKeys {
		db.getMesh(res, m)
		if res.VerticalLevel == y {
			return res
		} else {
			fmt.Println("mesh found but VerticalLevel failed", x, z, y, m)
		}
	}
	return nil
}
