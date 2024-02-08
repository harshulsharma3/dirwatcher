package api

import (
	"database/sql"
	"net/http"
	"time"

	"dirwatcher/config"
	"dirwatcher/database"
	"dirwatcher/dirwatcher"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(router *gin.Engine, db *sql.DB, startTask chan bool, stopTask chan bool) {
	var lastTaskID int64

	// Endpoint to configure a new task
	router.POST("/config", func(c *gin.Context) {
		var newTask config.Task
		if err := c.BindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration"})
			return
		}

		id, err := database.InsertTask(db, newTask)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert task into the database"})
			return
		}

		go dirwatcher.StartBackgroundTask(c, newTask, db, startTask, stopTask) // Start the background task
		startTask <- true

		lastTaskID = id
		c.Set("id", id)
		c.Set("start", time.Now().Format("2006-01-02 15:04:05"))
		newTask.TaskID = int(id)

		c.JSON(http.StatusOK, gin.H{"message": "Configuration updated", "id": id})
	})

	// Get task results endpoint
	router.GET("/results", func(c *gin.Context) {
		results, err := database.GetTaskResults(c, db, lastTaskID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task results"})
			return
		}
		c.JSON(http.StatusOK, results)
	})

	// Stop task endpoint
	router.POST("/stop", func(c *gin.Context) {
		stopTask <- true // Signal the background task to stop
		lastTaskID = 0
		c.JSON(http.StatusOK, gin.H{"message": "Task stopped"})
	})
}
