package dirwatcher

import (
	"database/sql"
	"dirwatcher/config"
	"dirwatcher/database"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
)

func StartBackgroundTask(c *gin.Context, newTask config.Task, db *sql.DB, startTask chan bool, stopTask chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	ids := c.Query("id")
	fmt.Println(ids)

	err = watcher.Add(newTask.Directory)
	if err != nil {
		fmt.Println("Error watching directory:", err)
		return
	}

	ticker := time.NewTicker(time.Duration(newTask.Interval) * time.Minute)
	defer ticker.Stop()

	running := false

	for {
		select {
		case event := <-watcher.Events:
			fmt.Println("watcher events :", event)
			if running {
				processFileEvent(c, newTask, event, db) // Handle file change immediately
			}
		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)
			database.UpdatefinalStatus(c, "Failed", db)
		case <-ticker.C:
			if running {
				processDirectory(c, newTask, db) // Run periodic task
			}
		case <-startTask:
			if !running {
				running = true
				processDirectory(c, newTask, db) // Start initial processing
			}
		case <-stopTask:
			running = false
			newTask = config.Task{}
			database.UpdatefinalStatus(c, "success", db)
			break
		}
	}
}

func processFileEvent(c *gin.Context, newTask config.Task, event fsnotify.Event, db *sql.DB) {
	filePath := event.Name

	// Handle file creation events
	if event.Op&fsnotify.Create == fsnotify.Create {
		fmt.Println("File created:", filePath)
		database.CountOccurrencesAndUpdateDB(c, newTask, filePath, db)
	}

	//Handle file write events
	if event.Op&fsnotify.Write == fsnotify.Write {
		fmt.Println("File writed :", filePath)
		database.CountOccurrencesAndUpdateDB(c, newTask, filePath, db)
	}

	// Handle file deletion events
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		fmt.Println("File deleted:", filePath)
		database.UpdateDBForFileDeletion(c, filePath, db)
	}

}

func processDirectory(c *gin.Context, newTask config.Task, db *sql.DB) {
	// ... (Read files in the directory, count occurrences, update database)
	files, err := ioutil.ReadDir(newTask.Directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	fmt.Println(files)
	for _, file := range files {
		filePath := filepath.Join(newTask.Directory, file.Name())
		database.CountOccurrencesAndUpdateDB(c, newTask, filePath, db)
	}
}
