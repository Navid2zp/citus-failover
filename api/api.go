package api

import (
	"github.com/Navid2zp/citus-failover/config"
	"github.com/Navid2zp/citus-failover/logging"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var logger *logging.Logger

func initAPI() {
	logger = logging.NewLogger("api")
	if !config.Config.Settings.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func adminMiddleware(c *gin.Context) {
	secret := c.Request.Header.Get("Secret")
	if secret == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrCodeSecretMissing})
		return
	}
	if secret != config.Config.API.Secret {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrCodeInvalidSecret})
		return
	}
	c.Next()
}

func Serve() {
	initAPI()

	router := gin.Default()
	v1 := router.Group("/v1", adminMiddleware)

	{
		v1.GET("/workers", getPrimaryWorkers)
		v1.GET("/coordinator", getPrimaryCoordinator)
		v1.GET("/coordinators", getAllCoordinators)
		v1.GET("/nodes", getNodes)
	}

	log.Println("Starting api on port", config.Config.API.Port, "...")

	err := router.Run(":" + config.Config.API.Port)
	if err != nil {
		logger.APIStartFailed(err)
	}
}
