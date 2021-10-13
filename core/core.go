package core

import (
	"database/sql"
	"errors"
	"time"
)

var ErrDBNotFound = errors.New("no such database found")

const nodeColumns = "formationid,nodeid,groupid,nodename,nodehost,nodeport,sysidentifier," +
	"goalstate,reportedstate,reportedpgisrunning,reportedrepstate,reporttime,reportedlsn," +
	"reportedtli,walreporttime,health,healthchecktime,statechangetime,candidatepriority," +
	"replicationquorum,nodecluster"

// Node represents a node in monitoring service.
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
	ReportedTLI         int       `db:"reportedtli" json:"reported_tli"`
	WALReportTime       time.Time `db:"walreporttime" json:"wal_report_time"`
	Health              int       `db:"health" json:"health"`
	HealthCheckTime     time.Time `db:"healthchecktime" json:"health_check_time"`
	StateChangeTime     time.Time `db:"statechangetime" json:"state_change_time"`
	CandidatePriority   int       `db:"candidatepriority" json:"candidate_priority"`
	ReplicationQuorum   bool      `db:"replicationquorum" json:"replication_quorum"`
	NodeCluster         string    `db:"nodecluster" json:"node_cluster"`
	IsCoordinator       bool      `db:"-" json:"-"`
}

const workerColumns = "nodeid,groupid,nodename,nodeport,noderack,hasmetadata,isactive,noderole," +
	"nodecluster,metadatasynced,shouldhaveshards"

// Worker represents a worker in coordinators.
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

const coordinatorColumns = "primary_node_id,primary_name,primary_host,primary_port"

// Coordinator represents a coordinator database
type Coordinator struct {
	PrimaryNodeID int    `db:"primary_node_id" json:"primary_node_id"`
	Name          string `db:"primary_name" json:"name"`
	Host          string `db:"primary_host" json:"host"`
	Port          int    `db:"primary_port" json:"port"`
}

// GetNodes returns all the node monitored in the monitor.
func GetNodes() ([]*Node, error) {
	var nodes []*Node
	err := monitorDB.Select(&nodes, "select "+nodeColumns+" from pgautofailover.node")
	return nodes, err
}

// GetPrimaryCoordinator returns the primary coordinator info given the database name.
func GetPrimaryCoordinator(dbname string) (*Coordinator, error) {
	var db *database
	if db = findDatabase(dbname); db == nil {
		return nil, ErrDBNotFound
	}
	return db.getCoordinator()
}

// GetPrimaryWorkers returns all the workers available in a database.
func GetPrimaryWorkers(dbname string) ([]*Worker, error) {
	var db *database
	var workers []*Worker
	if db = findDatabase(dbname); db == nil {
		return workers, nil
	}
	return db.getPrimaryWorkers()
}

// GetCoordinators returns a list of all the primary and non-primary coordinator nodes.
func GetCoordinators(dbname string) ([]*Node, error) {
	var coordinators []*Node
	var db *database
	if db = findDatabase(dbname); db == nil {
		return coordinators, nil
	}
	err := monitorDB.Select(
		&coordinators,
		"select "+nodeColumns+" from pgautofailover.node where formationid = $1;",
		db.formation,
	)
	return coordinators, err
}

// getPrimaryWorkers returns all the primary workers for a database.
func (d *database) getPrimaryWorkers() ([]*Worker, error) {
	var workers []*Worker
	err := d.db.Select(&workers,
		"SELECT "+workerColumns+" from pg_dist_node where noderole = 'primary';")
	return workers, err
}

// getCoordinator returns the coordinator for the database.
func (d *database) getCoordinator() (*Coordinator, error) {
	var node Coordinator
	err := monitorDB.Get(&node,
		"select "+coordinatorColumns+" from pgautofailover.get_primary($1);", d.formation)
	return &node, err
}

// isPrimary checks if the worker is a primary node in the monitor.
func (w *Worker) isPrimary() (bool, *Node, error) {
	var newNode Node
	// goalstate will be `wait_primary` when there is no second node as backup
	// or a node is not verified as secondary yet
	// not including `wait_primary` causes the primary check to fail when there is only one healthy node
	// https://pg-auto-failover.readthedocs.io/en/master/tutorial.html#cause-a-node-failure
	// https://github.com/Navid2zp/citus-failover/issues/1
	err := monitorDB.Get(&newNode, `select `+nodeColumns+` from pgautofailover.node
		where formationid = 
		      (select formationid from pgautofailover.node where nodehost = $1 and nodeport = $2 limit 1)
		  and (goalstate = 'primary' or goalstate = 'wait_primary')
		  and (nodehost != $1 or nodeport != $2) limit 1;`,
		w.Host, w.Port)
	if err == sql.ErrNoRows {
		return true, nil, nil
	}
	return false, &newNode, err
}

// updateWorkerInCoordinator updates a worker node in the database.
func (w *Worker) updateWorkerInCoordinator(newHost string, newPort int, db *database) error {
	_, err := db.db.Exec(`select * from citus_update_node($1, $2, $3);`, w.ID, newHost, newPort)
	return err
}

// connect connects to the database
// checks the coordinator state and the previous connections
// establishes a new one if connection is lost primary node changed
func (d *database) connect(coordinatorNode *Coordinator) error {
	var err error
	if d.db == nil {
		d.host = coordinatorNode.Host
		d.port = coordinatorNode.Port
		d.db, err = openDBConnection(d.host, d.username, d.dbname, d.password, d.port, d.sslmode)
		return err
	}
	if d.host != coordinatorNode.Host || d.port != coordinatorNode.Port {
		logger.CoordinatorChanged(coordinatorNode.Host, d.host, d.dbname, coordinatorNode.Port, d.port)
		d.host = coordinatorNode.Host
		d.port = coordinatorNode.Port
		d.db, err = openDBConnection(d.host, d.username, d.dbname, d.password, d.port, d.sslmode)
		return err
	}
	if d.db.Ping() != nil {
		logger.CoordinatorConnectionLost(d.host, d.dbname, d.username, d.port)
		d.db, err = openDBConnection(d.host, d.username, d.dbname, d.password, d.port, d.sslmode)
		return err
	}
	return nil
}
