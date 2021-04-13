package core

import (
	"database/sql"
	"github.com/Navid2zp/citus-failover/config"
	"time"
)

func validateWorker(worker *Worker) {
	isPrimary, newNode, err := worker.isPrimary()
	if err != nil {
		logger.WorkerPrimaryCheckFailed(err, worker.ID)
		return
	}
	if !isPrimary {
		logger.WorkerStateChange(worker.Host, newNode.Host, worker.Port, newNode.Port)
		err = worker.updateCoordinator(newNode.Host, newNode.Port)
		if err != nil {
			logger.WorkerUpdateFailed(err)
			return
		}
		logger.WorkerUpdated(worker.Host, newNode.Host, worker.Port, newNode.Port)
	}
}

func validateWorkers() {
	workers, err := GetPrimaryWorkers()
	if err != nil {
		logger.GetWorkersFailed(err)
		return
	}
	for _, worker := range workers {
		go validateWorker(worker)
	}
}

func Monitor() {
	for {
		time.Sleep(time.Duration(config.Config.Settings.CheckInterval) * time.Millisecond)
		coordinator, err := GetCoordinator()
		if err != nil {
			if err == sql.ErrNoRows {
				logger.NoCoordinatorFound()
			} else {
				logger.GetCoordinatorFailed(err)
			}
			continue
		}
		err = coordinator.connect()
		if err != nil {
			logger.CoordinatorConnectionFailed(coordinator.Host, coordinator.Port, err)
			continue
		}
		validateWorkers()
	}
}
