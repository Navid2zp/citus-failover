package api

import (
	"github.com/Navid2zp/citus-failover/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrCodeServerError   = "SERVER_ERROR"
	ErrCodeSecretMissing = "SECRET_MISSING"
	ErrCodeInvalidSecret = "INVALID_SECRET"
)

func getPrimaryWorkers(c *gin.Context) {
	workers, err := core.GetPrimaryWorkers()
	if err != nil {
		logger.GetWorkersFailed(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   ErrCodeServerError,
			"message": "failed to get workers from coordinator",
		})
		return
	}

	c.JSON(http.StatusOK, workers)
}

func getPrimaryCoordinator(c *gin.Context) {
	coordinator, err := core.GetCoordinator()
	if err != nil {
		logger.GetCoordinatorFailed(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   ErrCodeServerError,
			"message": "failed to get coordinator from monitor",
		})
		return
	}

	c.JSON(http.StatusOK, coordinator)
}

func getAllCoordinators(c *gin.Context) {
	coordinators, err := core.GetAllCoordinators()
	if err != nil {
		logger.GetCoordinatorFailed(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   ErrCodeServerError,
			"message": "failed to get coordinators from monitor",
		})
		return
	}

	c.JSON(http.StatusOK, coordinators)
}

func getNodes(c *gin.Context) {
	nodes, err := core.GetNodes()
	if err != nil {
		logger.GetWorkersFailed(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   ErrCodeServerError,
			"message": "failed to get nodes from monitor",
		})
		return
	}

	c.JSON(http.StatusOK, nodes)
}
