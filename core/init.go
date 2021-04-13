package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Navid2zp/citus-failover/config"
	"github.com/Navid2zp/citus-failover/logging"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

var monitorDB *sqlx.DB
var logger *logging.Logger
var currentCoordinator struct {
	host string
	port int
	db   *sqlx.DB
}

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

func openCoordinatorConnection(host string, port int) error {
	var err error
	currentCoordinator.host = ""
	currentCoordinator.port = 0
	currentCoordinator.db, err = sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			config.Config.Coordinator.Username,
			config.Config.Coordinator.Password,
			config.Config.Coordinator.DBName,
		),
	)
	if err == nil {
		currentCoordinator.host = host
		currentCoordinator.port = port
	}
	return err
}

func InitMonitor() {
	openMonitoringConnection()
	logger = logging.NewLogger("core")
}

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) ToInt64() int64 {
	return ni.Int64
}

func Int64ToNullInt64(i int64) NullInt64 {
	if i == 0 {
		return NullInt64{}
	}
	return NullInt64{
		sql.NullInt64{
			Int64: i,
			Valid: true,
		},
	}
}

func TimeToNullTime(t time.Time, null bool) NullTime {
	if null {
		return NullTime{}
	}
	return NullTime{
		pq.NullTime{
			Time:  t,
			Valid: true,
		},
	}
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = err == nil
	return err
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for NullBool
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nb.Bool)
	nb.Valid = err == nil
	return err
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

//UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = err == nil
	return err
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

func (ns *NullString) ToString() string {
	return ns.String
}

func StringToNullString(s string) NullString {
	if s == "" {
		return NullString{}
	}
	return NullString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

func Float64ToNullFloat64(f float64) NullFloat64 {
	if f == 0 {
		return NullFloat64{}
	}
	return NullFloat64{
		sql.NullFloat64{
			Float64: f,
			Valid:   true,
		},
	}
}

type NullTime struct {
	pq.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nt.Time)
	nt.Valid = err == nil
	return err

}
