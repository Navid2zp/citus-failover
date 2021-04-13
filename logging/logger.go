package logging

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type Logger struct {
	service string
	l       *zap.Logger
}

func NewLogger(service string) *Logger {
	//zap.NewProductionConfig()
	//cfg := zap.NewProductionConfig()
	//cfg.OutputPaths = []string{"stdout"}
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return &Logger{service, logger}
}

func (l *Logger) MonitorInit() {
	l.l.Info("monitoring initialized ...")
}

func (l *Logger) NoCoordinatorFound() {
	l.l.Error("no coordinator found! probably switching state, skipping ...",
		zap.String("service", l.service))
}

func (l *Logger) CoordinatorChanged(oldHost, newHost string, oldPort, newPort int) {
	l.l.Info("coordinator changed, connecting to new coordinator.",
		zap.String("service", l.service),
		zap.String("old-coordinator", fmt.Sprintf("%s:%d", oldHost, oldPort)),
		zap.String("new-coordinator", fmt.Sprintf("%s:%d", newHost, newPort)))
}

func (l *Logger) CoordinatorConnectionLost(host string, port int) {
	l.l.Warn("coordinator connection lost, reconnecting ...",
		zap.String("service", l.service),
		zap.String("host", host), zap.Int("port", port))
}

func (l *Logger) CoordinatorConnectionFailed(host string, port int, err error) {
	l.l.Error("connecting to coordinator failed!",
		zap.String("service", l.service),
		zap.String("host", host), zap.Int("port", port), zap.Error(err))
}

func (l *Logger) GetCoordinatorFailed(err error) {
	l.l.Error("failed to get coordinator, skipping ...",
		zap.String("service", l.service), zap.Error(err))
}

func (l *Logger) GetWorkersFailed(err error) {
	if err == sql.ErrNoRows {
		l.NoWorkersFound()
		return
	}
	l.l.Error("failed to get workers!",
		zap.String("service", l.service), zap.Error(err))
}

func (l *Logger) WorkerPrimaryCheckFailed(err error, workerID int) {
	l.l.Error("failed to check if worker is primary!",
		zap.String("service", l.service), zap.Error(err), zap.Int("worker-id", workerID))
}

func (l *Logger) WorkerStateChange(oldHost, newHost string, oldPort, newPort int) {
	oldWorker := fmt.Sprintf("%s:%d", oldHost, oldPort)
	newWorker := fmt.Sprintf("%s:%d", newHost, newPort)
	l.l.Info("worker change detected, updating coordinator ...",
		zap.String("service", l.service),
		zap.String("old-worker", oldWorker), zap.String("new-worker", newWorker))
}

func (l *Logger) WorkerUpdated(oldHost, newHost string, oldPort, newPort int) {
	oldWorker := fmt.Sprintf("%s:%d", oldHost, oldPort)
	newWorker := fmt.Sprintf("%s:%d", newHost, newPort)
	l.l.Info("worker updated",
		zap.String("service", l.service),
		zap.String("old-worker", oldWorker), zap.String("new-worker", newWorker))
}

func (l *Logger) WorkerUpdateFailed(err error) {
	l.l.Error("failed to update worker in coordinator",
		zap.String("service", l.service), zap.Error(err))
}

func (l *Logger) NoWorkersFound() {
	l.l.Error("no workers found!",
		zap.String("service", l.service))
}
