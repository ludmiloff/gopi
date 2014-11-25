package gopi

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"log"
	"strings"
)

func (this *Application) InitDB() {
	if this.Config.Has("mysql") {
		if this.Config.Has("mysql.simple") {
			this.DB = GetDBSimple(
				this.Config.Get("dsn").(string),
				this.Config.Get("engine").(string),
				this.Config.Get("encoding").(string))

		} else { // TODO: complex mysql setup

		}
	}
}

func GetDB(user, password, hostname, port, database, encoding, engine string) *gorp.DbMap {
	db, err := sql.Open("mysql", fmt.Sprint(user, ":", password, "@(", hostname, ":", port, ")/", database, "?charset=", encoding))
	checkErr(err, "sql.Open failed")

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{engine, strings.ToUpper(encoding)}}

	return dbMap
}

func GetDBSimple(dsn, engine, encoding string) *gorp.DbMap {
	db, err := sql.Open("mysql", dsn)
	checkErr(err, "sql.Open failed")

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{engine, strings.ToUpper(encoding)}}

	return dbMap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
