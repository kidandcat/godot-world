package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/tidwall/buntdb"
)

type DB struct {
	db *buntdb.DB
}

func dbInitialize() *DB {
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

func dbNextIndex(tx *buntdb.Tx) string {
	index, e := tx.Get("index")
	utilsCheck(e)
	iindex, e := strconv.ParseInt(index, 10, 64)
	utilsCheck(e)
	next := iindex + 1
	nextIndexStr := fmt.Sprintf("%v", next)
	tx.Set("index", nextIndexStr, nil)
	return nextIndexStr
}
