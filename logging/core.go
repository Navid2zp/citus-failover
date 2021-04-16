package logging

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

func (l *Logger) MonitorInit() {
	l.l.Info("monitoring initialized ...")
}

func (l *Logger) NoCoordinatorFound() {
	l.l.Error("no coordinator found! probably switching state, skipping ...",
		zap.String("service", l.service))
}

func (l *Logger) CoordinatorChanged(oldHost, newHost, dbname string, oldPort, newPort int) {
	l.l.Info("coordinator changed, connecting to new coordinator.",
		zap.String("service", l.service), zap.String("dbname", dbname),
		zap.String("old-coordinator", fmt.Sprintf("%s:%d", oldHost, oldPort)),
		zap.String("new-coordinator", fmt.Sprintf("%s:%d", newHost, newPort)))
}

func (l *Logger) CoordinatorConnectionLost(host, dbname, dbUsername string, port int) {
	l.l.Warn("coordinator connection lost, reconnecting ...",
		zap.String("service", l.service),
		zap.String("host", host), zap.String("database", dbname), zap.String("username", dbUsername),
		zap.Int("port", port))
}

func (l *Logger) CoordinatorConnectionFailed(host, dbname, dbUsername string, port int, err error) {
	l.l.Error("connecting to coordinator failed!",
		zap.String("service", l.service),
		zap.String("host", host), zap.String("database", dbname), zap.String("username", dbUsername),
		zap.Int("port", port), zap.Error(err))
}

func (l *Logger) GetCoordinatorFailed(err error) {
	l.l.Error("failed to get coordinator, skipping ...",
		zap.String("service", l.service), zap.Error(err))
}

func (l *Logger) GetWorkersFailed(err error, dbname string) {
	if err == sql.ErrNoRows {
		l.NoWorkersFound(dbname)
		return
	}
	l.l.Error("failed to get workers!",
		zap.String("service", l.service), zap.String("dbname", dbname), zap.Error(err))
}

func (l *Logger) WorkerPrimaryCheckFailed(err error, workerID int, dbname string) {
	l.l.Error("failed to check if worker is primary!",
		zap.String("service", l.service), zap.String("database", dbname),
		zap.Error(err), zap.Int("worker-id", workerID))
}

func (l *Logger) WorkerStateChange(oldHost, newHost, dbname string, oldPort, newPort int) {
	oldWorker := fmt.Sprintf("%s:%d", oldHost, oldPort)
	newWorker := fmt.Sprintf("%s:%d", newHost, newPort)
	l.l.Info("worker change detected, updating coordinator ...",
		zap.String("service", l.service), zap.String("database", dbname),
		zap.String("old-worker", oldWorker), zap.String("new-worker", newWorker))
}

func (l *Logger) WorkerUpdated(oldHost, newHost, dbname string, oldPort, newPort int) {
	oldWorker := fmt.Sprintf("%s:%d", oldHost, oldPort)
	newWorker := fmt.Sprintf("%s:%d", newHost, newPort)
	l.l.Info("worker updated",
		zap.String("service", l.service), zap.String("database", dbname),
		zap.String("old-worker", oldWorker), zap.String("new-worker", newWorker))
}

func (l *Logger) WorkerUpdateFailed(err error, dbname string) {
	l.l.Error("failed to update worker in coordinator",
		zap.String("service", l.service), zap.String("dbname", dbname), zap.Error(err))
}

func (l *Logger) NoWorkersFound(dbname string) {
	l.l.Error("no workers found!",
		zap.String("service", l.service), zap.String("dbname", dbname))
}
