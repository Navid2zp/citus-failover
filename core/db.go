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
	sslmode   string `default:"disable"`
	db        *sqlx.DB
}

// databases holds the list of all databases to be monitored
var databases []*database

// monitorDB is database instance for monitor
var monitorDB *sqlx.DB

// openMonitoringConnection opens a connection to monitor database
func openMonitoringConnection() {
	var err error

	monitorDB, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Config.Monitor.Host,
			config.Config.Monitor.Port,
			config.Config.Monitor.User,
			config.Config.Monitor.Password,
			config.Config.Monitor.DBName,
			config.Config.Monitor.SSLMode,
		),
	)
	if err != nil {
		panic(err)
	}
}

// setupDatabases adds all the listed databases in the config file to databases.
func setupDatabases() {
	for _, db := range config.Config.Coordinators {
		coordinator := database{
			formation: db.Formation,
			host:      "",
			port:      0,
			username:  db.Username,
			password:  db.Password,
			dbname:    db.DBName,
			sslmode:   db.SSLMode,
			db:        nil,
		}
		databases = append(databases, &coordinator)
	}
}

// openDBConnection opens a database connection.
func openDBConnection(host, username, dbname, password string, port int, sslmode string) (*sqlx.DB, error) {
	return sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host,
			port,
			username,
			password,
			dbname,
			sslmode,
		),
	)
}

// findDatabase finds a database in monitoring databases list given the database name.
func findDatabase(dbname string) *database {
	for _, db := range databases {
		if db.dbname == dbname {
			return db
		}
	}
	return nil
}
