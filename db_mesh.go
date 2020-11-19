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
			kID := strings.Split(key, ":")[1]
			vl, err := tx.Get("mesh:" + kID + ":verticalLevel")
			if err == nil {
				verticalLevel, err := strconv.ParseInt(vl, 10, 64)
				utilsCheck(err)
				if verticalLevel != p.VerticalLevel { // do not delete if it is not in the same vertical level
					return true
				}
			}
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
		tx.Set("mesh:"+index+":rotation", p.Rotation, nil)
		p.ID = index
		return nil
	})
	utilsCheck(err)
}

func (db *DB) getMesh(p *tMesh, ID string) {
	db.db.View(func(tx *buntdb.Tx) error {
		rot, err := tx.Get("mesh:" + ID + ":rotation")
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

		if rot == "" {
			rot = "down" // down default value
		}

		var X, Z int64
		fmt.Sscanf(pos, "[%d %d]", &X, &Z)
		p.ID = ID
		p.Type, err = strconv.ParseInt(T, 10, 64)
		p.Position.X = X
		p.Position.Z = Z
		p.Rotation = rot
		p.VerticalLevel, err = strconv.ParseInt(verticalLevel, 10, 64)
		p.Walkable = utilsStringToBool(walkable)
		p.WalkingCost, err = strconv.Atoi(walkingCost)
		utilsCheck(err)
		return nil
	})
}

func (db *DB) deleteMesh(ID string) {
	db.db.View(func(tx *buntdb.Tx) error {
		tx.Delete("mesh:" + ID + ":rotation")
		tx.Delete("mesh:" + ID + ":type")
		tx.Delete("mesh:" + ID + ":pos")
		tx.Delete("mesh:" + ID + ":verticalLevel")
		tx.Delete("mesh:" + ID + ":walkable")
		tx.Delete("mesh:" + ID + ":walkingCost")
		return nil
	})
}

func (db *DB) deleteMeshByPos(x, z, y int64) {
	meshAux := &tMesh{}
	meshesKeys := []string{}
	_x := strconv.FormatInt(x, 10)
	_z := strconv.FormatInt(z, 10)
	db.db.View(func(tx *buntdb.Tx) error {
		return tx.Intersects("meshPos", "["+_x+" "+_z+"],["+_x+" "+_z+"]", func(key, val string) bool {
			kID := strings.Split(key, ":")[1]
			meshesKeys = append(meshesKeys, kID)
			db.getMesh(meshAux, kID)
			if meshAux.VerticalLevel == y {
				db.deleteMesh(kID)
			}
			return true
		})
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
		}
	}
	return nil
}
