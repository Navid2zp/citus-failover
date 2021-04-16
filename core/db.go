package core

import (
	"fmt"
	"github.com/Navid2zp/citus-failover/config"
	"github.com/jmoiron/sqlx"
)

type database struct {
	formation string
	host      string
	port      int
	username  string
	password  string
	dbname    string
	db        *sqlx.DB
}

var databases []*database
var monitorDB *sqlx.DB

func openMonitoringConnection() {
	var err error

	monitorDB, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Config.Monitor.Host,
			config.Config.Monitor.Port,
			config.Config.Monitor.User,
			config.Config.Monitor.Password,
			config.Config.Monitor.DBName,
		),
	)
	if err != nil {
		panic(err)
	}
}

func setupDatabases() {
	for _, db := range config.Config.Coordinators {
		coordinator := database{
			formation: db.Formation,
			host:      "",
			port:      0,
			username:  db.Username,
			password:  db.Password,
			dbname:    db.DBName,
			db:        nil,
		}
		databases = append(databases, &coordinator)
	}
}

func openDBConnection(host, username, dbname, password string, port int) (*sqlx.DB, error) {
	return sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			username,
			password,
			dbname,
		),
	)
}

func findDatabase(dbname string) *database {
	for _, db := range databases {
		if db.dbname == dbname {
			return db
		}
	}
	return nil
}
