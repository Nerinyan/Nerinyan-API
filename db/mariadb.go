package db

import (
	"Nerinyan-API/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pterm/pterm"
)

var Maria *sql.DB

func ConnectMaria() {
	s := &config.Config.Sql
	db, err := sql.Open("mysql", s.Id+":"+s.Passwd+"@tcp("+s.Url+")/")
	if Maria = db; db != nil {
		Maria.SetMaxOpenConns(100)
		if _, err = Maria.Exec("SET SQL_SAFE_UPDATES = 0;"); err != nil {
			panic(err)
		}
		pterm.Success.Println("RDBMS Connected.")
	} else {
		panic(err)
	}
}
