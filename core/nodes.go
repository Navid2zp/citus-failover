package core

import "time"

type Node struct {
	FormationID         string
	ID                  int
	GroupID             int
	NodeName            string
	Host                string
	Port                int
	SysIdentifier       string
	GoalState           string
	ReportedState       string
	ReportedPGIsRunning bool
	ReportedRepState    string
	ReportTime          time.Time
	ReportedLSN         string
	WALReportTime       time.Time
	Health              int
	HealthCheckTime     time.Time
	StateChangeTime     time.Time
	CandidatePriority   int
	ReplicationQuorum   bool
	NodeCluster         string
}
