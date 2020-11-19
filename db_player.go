package main

import (
	"fmt"

	"github.com/tidwall/buntdb"
)

/* ------ */
/* PLAYER */
/* ------ */

func (db *DB) newPlayer(p *tPlayer) {
	if p == nil {
		panic("newPlayer: Player is NIL")
	}
	err := db.db.Update(func(tx *buntdb.Tx) error {
		index := dbNextIndex(tx)
		tx.Set("player:"+index+":name", p.Name, nil)
		tx.Set("player:"+index+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		p.ID = index
		return nil
	})
	utilsCheck(err)
}

func (db *DB) setPlayerName(p *tPlayer) {
	err := db.db.Update(func(tx *buntdb.Tx) error {
		tx.Set("player:"+p.ID+":name", p.Name, nil)
		return nil
	})
	utilsCheck(err)
}

func (db *DB) setPlayerPos(p *tPlayer) {
	err := db.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("player:"+p.ID+":pos", fmt.Sprintf("[%v %v]", p.Position.X, p.Position.Z), nil)
		return err
	})
	utilsCheck(err)
}

func (db *DB) getPlayerName(p *tPlayer) {
	db.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("player:" + p.ID + ":name")
		if err != nil {
			return err
		}
		p.Name = val
		return nil
	})
}

func (db *DB) getPlayerPos(p *tPlayer) {
	err := db.db.View(func(tx *buntdb.Tx) error {
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
	utilsCheck(err)
}

func (db *DB) getPlayer(ID string) *tPlayer {
	p := &tPlayer{}
	err := db.db.View(func(tx *buntdb.Tx) error {
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
	utilsCheck(err)
	return p
}
