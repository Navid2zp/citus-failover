package core

import (
	"database/sql"
	"github.com/Navid2zp/citus-failover/config"
	"time"
)

func (d *database) validateWorker(worker *Worker) {
	isPrimary, newNode, err := worker.isPrimary()
	if err != nil {
		logger.WorkerPrimaryCheckFailed(err, worker.ID, d.dbname)
		return
	}
	if !isPrimary {
		logger.WorkerStateChange(worker.Host, newNode.Host, d.dbname, worker.Port, newNode.Port)
		err = worker.updateCoordinator(newNode.Host, newNode.Port, d)
		if err != nil {
			logger.WorkerUpdateFailed(err, d.dbname)
			return
		}
		logger.WorkerUpdated(worker.Host, newNode.Host, d.dbname, worker.Port, newNode.Port)
	}
}

func (d *database) validateWorkers() {
	workers, err := d.getPrimaryWorkers()
	if err != nil {
		logger.GetWorkersFailed(err, d.dbname)
		return
	}
	for _, worker := range workers {
		go d.validateWorker(worker)
	}
}

func (d *database) monitor() {
	for {
		time.Sleep(time.Duration(config.Config.Settings.CheckInterval) * time.Millisecond)
		coordinator, err := d.getCoordinator()
		if err != nil {
			if err == sql.ErrNoRows {
				logger.NoCoordinatorFound()
			} else {
				logger.GetCoordinatorFailed(err)
			}
			continue
		}
		err = d.connect(coordinator)
		if err != nil {
			logger.CoordinatorConnectionFailed(
				coordinator.Host, d.dbname, d.username, coordinator.Port, err)
			continue
		}
		d.validateWorkers()
	}
}


func Monitor() {

	for _, db := range databases {
		go db.monitor()
	}
	select {}
}
