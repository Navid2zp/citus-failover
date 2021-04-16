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

	// ErrCodeBadData represents the code for providing bad or wrong data by user in an API request.
	ErrCodeBadData = "BAD_DATA"
)

// getPrimaryWorkers returns a list of all primary worker nodes for the given database.
func getPrimaryWorkers(c *gin.Context) {
	databaseName := c.Param("database")
	workers, err := core.GetPrimaryWorkers(databaseName)

	if err != nil {
		logger.GetWorkersFailed(err, "")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   ErrCodeServerError,
			"message": "failed to get workers from coordinator",
		})
		return
	}

	c.JSON(http.StatusOK, workers)
}

// getPrimaryCoordinator returns the primary coordinator node in the monitor.
func getPrimaryCoordinator(c *gin.Context) {
	databaseName := c.Param("database")
	coordinator, err := core.GetPrimaryCoordinator(databaseName)
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

// getAllCoordinators returns a list of all coordinator nodes available in the monitor.
func getAllCoordinators(c *gin.Context) {
	databaseName := c.Param("database")
	coordinators, err := core.GetCoordinators(databaseName)
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

// getNodes returns a list of all available nodes in the monitor.
func getNodes(c *gin.Context) {
	nodes, err := core.GetNodes()
	if err != nil {
		logger.GetWorkersFailed(err, "")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   ErrCodeServerError,
			"message": "failed to get nodes from monitor",
		})
		return
	}

	c.JSON(http.StatusOK, nodes)
}

