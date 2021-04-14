package core

import (
	"database/sql"
	"github.com/Navid2zp/citus-failover/config"
	"time"
)

type Node struct {
	FormationID         string    `db:"formationid" json:"formation_id"`
	ID                  int       `db:"nodeid" json:"id"`
	GroupID             int       `db:"groupid" json:"group_id"`
	Name                string    `db:"nodename" json:"name"`
	Host                string    `db:"nodehost" json:"host"`
	Port                int       `db:"nodeport" json:"port"`
	SysIdentifier       string    `db:"sysidentifier" json:"sys_identifier"`
	GoalState           string    `db:"goalstate" json:"goal_state"`
	ReportedState       string    `db:"reportedstate" json:"reported_state"`
	ReportedPGIsRunning bool      `db:"reportedpgisrunning" json:"reported_pg_is_running"`
	ReportedRepState    string    `db:"reportedrepstate" json:"reported_rep_state"`
	ReportTime          time.Time `db:"reporttime" json:"report_time"`
	ReportedLSN         string    `db:"reportedlsn" json:"reported_lsn"`
	WALReportTime       time.Time `db:"walreporttime" json:"wal_report_time"`
	Health              int       `db:"health" json:"health"`
	HealthCheckTime     time.Time `db:"healthchecktime" json:"health_check_time"`
	StateChangeTime     time.Time `db:"statechangetime" json:"state_change_time"`
	CandidatePriority   int       `db:"candidatepriority" json:"candidate_priority"`
	ReplicationQuorum   bool      `db:"replicationquorum" json:"replication_quorum"`
	NodeCluster         string    `db:"nodecluster" json:"node_cluster"`
	IsCoordinator       bool      `db:"-" json:"-"`
}

type Worker struct {
	ID               int      `db:"nodeid" json:"id"`
	GroupID          int      `db:"groupid" json:"group_id"`
	Host             string   `db:"nodename" json:"host"`
	Port             int      `db:"nodeport" json:"port"`
	Rack             string   `db:"noderack" json:"rack"`
	HasMetaData      NullBool `db:"hasmetadata" json:"has_meta_data"`
	Active           bool     `db:"isactive" json:"active"`
	Role             string   `db:"noderole" json:"role"`
	Cluster          string   `db:"nodecluster" json:"cluster"`
	MetaDataSynced   NullBool `db:"metadatasynced" json:"meta_data_synced"`
	ShouldHaveShards bool     `db:"shouldhaveshards" json:"should_have_shards"`
}

type Coordinator struct {
	PrimaryNodeID int    `db:"primary_node_id" json:"primary_node_id"`
	Name          string `db:"primary_name" json:"name"`
	Host          string `db:"primary_host" json:"host"`
	Port          int    `db:"primary_port" json:"port"`
}

func GetNodes() ([]*Node, error) {
	var nodes []*Node
	err := monitorDB.Select(&nodes, `select * from pgautofailover.node`)
	return nodes, err
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

func GetAllCoordinators() ([]*Node, error) {
	var coordinators []*Node
	err := monitorDB.Select(&coordinators,
		`select * from pgautofailover.node where formationid = $1;`, config.Config.Coordinator.Formation)
	return coordinators, err
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
