package gopi

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pelletier/go-toml"

	"log"
	"strings"
)

func (this *Application) InitDB() {
	if this.Config.Has("mysql") {
		var mysql = this.Config.Get("mysql").(*toml.TomlTree)
		if mysql.Has("simple") {
			this.DB = GetMySQLDBSimple(
				mysql.Get("dsn").(string),
				mysql.Get("engine").(string),
				mysql.Get("encoding").(string))

		} else { // TODO: complex mysql setup

		}
	}
}

func GetMySQLDB(user, password, hostname, port, database, encoding, engine string) *gorp.DbMap {
	db, err := sql.Open("mysql", fmt.Sprint(user, ":", password, "@(", hostname, ":", port, ")/", database, "?charset=", encoding))
	CheckErr(err, "sql.Open failed")

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{engine, strings.ToUpper(encoding)}}

	return dbMap
}

func GetMySQLDBSimple(dsn, engine, encoding string) *gorp.DbMap {
	db, err := sql.Open("mysql", dsn)
	CheckErr(err, "sql.Open failed")

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{engine, strings.ToUpper(encoding)}}

	log.Println("MYSQL configured")
	return dbMap
}

func CheckErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
