package database

import (
	"database/sql"
	"dirwatcher/config"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitDB(host, port, user, password string) (sportsDb *sql.DB, err error) {
	Db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/watcher")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = Db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return Db, nil
}

func GetTaskResults(c *gin.Context, db *sql.DB, ids int64) ([]config.TaskResult, error) {
	// id := c.GetInt64("id")
	// ids := c.Query("id")
	// fmt.Println(ids)
	rows, err := db.Query("SELECT d.dir_path,t.start_time, t.end_time, t.total_runtime, t.files_added, t.files_deleted,t.magic_string_count, t.status FROM directory as d left join tasks as t on d.id = t.dir_id where d.id = ?", ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []config.TaskResult
	for rows.Next() {
		var result config.TaskResult
		var add, del sql.NullString
		var count sql.NullInt64
		if err := rows.Scan(
			&result.Directory,
			&result.StartTime,
			&result.EndTime,
			&result.TotalRuntime,
			&add,
			&del,
			&count,
			&result.Status,
		); err != nil {
			return nil, err
		}
		if add.Valid {
			result.FilesAdded = add.String
		}
		if del.Valid {
			result.FilesDeleted = del.String
		}
		if count.Valid {
			result.MagicStringCount = strconv.Itoa(int(count.Int64))
		}
		results = append(results, result)
	}
	return results, nil
}

func InsertTask(db *sql.DB, newTask config.Task) (int64, error) {

	sqlstmt := `Insert into directory set dir_path =?, magic_string =? , time_interval =? , added_at = ?, status = ? `

	stmt, err := db.Prepare(sqlstmt)
	if err != nil {
		fmt.Println("Prepare", err.Error())
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newTask.Directory, newTask.MagicString, newTask.Interval, time.Now().Format("2006-01-02 15:04:05"), "in_progress")

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil

}

func UpdateDBForFileDeletion(c *gin.Context, filePath string, db *sql.DB) {
	fmt.Println("file deleted...", filePath)

	stime, _ := time.Parse("", c.GetString("start"))
	diff := time.Now().Sub(stime).Hours()

	_, err := db.Query("INSERT INTO tasks (dir_id, start_time, end_time, total_runtime, files_deleted, status) VALUES (?,?,?,?,?,?)", c.GetInt64("id"), c.GetString("start"), time.Now().Format("2006-01-02 15:04:05"), diff, filePath, "success")
	if err != nil {
		fmt.Println("error while inserting tasks after deleting file - ", err)
	}

}

func CountOccurrencesAndUpdateDB(c *gin.Context, newTask config.Task, filePath string, db *sql.DB) {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	occurrenceCount := strings.Count(string(contents), newTask.MagicString)

	stime, _ := time.Parse("", c.GetString("start"))
	diff := time.Now().Sub(stime).Hours()

	_, err = db.Query("INSERT INTO tasks (dir_id, start_time, end_time, total_runtime, files_added,  magic_string_count, status) VALUES (?,?,?,?,?,?,?)", c.GetInt64("id"), c.GetString("start"), time.Now().Format("2006-01-02 15:04:05"), diff, filePath, occurrenceCount, "success")
	if err != nil {
		fmt.Println("error while inserting tasks - ", err)
	}
}

func UpdatefinalStatus(c *gin.Context, status string, db *sql.DB) error {

	sqlstmt := `update directory set status =? where id = ? `

	stmt, err := db.Prepare(sqlstmt)
	if err != nil {
		fmt.Println("Prepare", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, c.GetInt64("id"))

	if err != nil {
		return err
	}

	return nil
}
