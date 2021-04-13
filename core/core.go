package core

import (
	"database/sql"
	"github.com/Navid2zp/citus-failover/config"
	"time"
)

type Node struct {
	FormationID         string    `db:"formationid"`
	ID                  int       `db:"nodeid"`
	GroupID             int       `db:"groupid"`
	Name                string    `db:"nodename"`
	Host                string    `db:"nodehost"`
	Port                int       `db:"nodeport"`
	SysIdentifier       string    `db:"sysidentifier"`
	GoalState           string    `db:"goalstate"`
	ReportedState       string    `db:"reportedstate"`
	ReportedPGIsRunning bool      `db:"reportedpgisrunning"`
	ReportedRepState    string    `db:"reportedrepstate"`
	ReportTime          time.Time `db:"reporttime"`
	ReportedLSN         string    `db:"reportedlsn"`
	WALReportTime       time.Time `db:"walreporttime"`
	Health              int       `db:"health"`
	HealthCheckTime     time.Time `db:"healthchecktime"`
	StateChangeTime     time.Time `db:"statechangetime"`
	CandidatePriority   int       `db:"candidatepriority"`
	ReplicationQuorum   bool      `db:"replicationquorum"`
	NodeCluster         string    `db:"nodecluster"`
	IsCoordinator       bool      `db:"-"`
}

type Worker struct {
	ID               int      `db:"nodeid"`
	GroupID          int      `db:"groupid"`
	Host             string   `db:"nodename"`
	Port             int      `db:"nodeport"`
	Rack             string   `db:"noderack"`
	HasMetaData      NullBool `db:"hasmetadata"`
	Active           bool     `db:"isactive"`
	Role             string   `db:"noderole"`
	Cluster          string   `db:"nodecluster"`
	MetaDataSynced   NullBool `db:"metadatasynced"`
	ShouldHaveShards bool     `db:"shouldhaveshards"`
}

type Coordinator struct {
	PrimaryNodeID int    `db:"primary_node_id"`
	Name          string `db:"primary_name"`
	Host          string `db:"primary_host"`
	Port          int    `db:"primary_port"`
}

func GetPrimaryWorkers() ([]*Worker, error) {
	var workers []*Worker
	err := currentCoordinator.db.Select(&workers, `SELECT * from pg_dist_node where noderole = 'primary';`)
	return workers, err
}

func (w *Worker) isPrimary() (bool, *Node, error) {
	var newNode Node
	err := monitorDB.Get(&newNode, `select * from pgautofailover.node
		where formationid = 
		      (select formationid from pgautofailover.node where nodehost = $1 and nodeport = $2 limit 1)
		  and goalstate = 'primary'
		  and (nodehost != $1 or nodeport != $2) limit 1;`,
		w.Host, w.Port)
	if err == sql.ErrNoRows {
		return true, nil, nil
	}
	return false, &newNode, err
}

func (w *Worker) updateCoordinator(newHost string, newPort int) error {
	_, err := currentCoordinator.db.Exec(`select * from citus_update_node($1, $2, $3);`, w.ID, newHost, newPort)
	return err
}

func GetCoordinator() (*Coordinator, error) {
	var node Coordinator
	err := monitorDB.Get(&node,
		`select * from pgautofailover.get_primary($1);`, config.Config.Coordinator.Formation)
	return &node, err
}

func (c *Coordinator) connect() error {

	if currentCoordinator.db == nil {
		return openCoordinatorConnection(c.Host, c.Port)
	}
	if c.Host != currentCoordinator.host || c.Port != currentCoordinator.port {
		logger.CoordinatorChanged(currentCoordinator.host, c.Host, currentCoordinator.port, c.Port)
		return openCoordinatorConnection(c.Host, c.Port)
	}
	if currentCoordinator.db.Ping() != nil {
		logger.CoordinatorConnectionLost(c.Host, c.Port)
		return openCoordinatorConnection(c.Host, c.Port)
	}
	return nil
}
