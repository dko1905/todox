package routes

import (
	"fmt"
	"hash/crc32"
	"net/http"
	"time"

	"0chaos.eu/todox/db"
	"0chaos.eu/todox/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Return all tasks.
// Cache using crc32 ETag.
func GetTasks(c *gin.Context) {
	conn, err := db.GetConnection(c)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer conn.Close()

	rows, err := conn.QueryContext(c, "SELECT id, title, done, created_ts FROM tasks;")
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer rows.Close()

	// Fetch tasks and calcaulate ETag
	tasks := []models.Task{}
	hash := crc32.NewIEEE()
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done, &task.CreatedTS)
		if err != nil {
			InternalServerError(c, err)
			return
		}

		// Append
		task.HashWrite32(&hash)
		tasks = append(tasks, task)
	}

	// Compare client and server ETags
	etag := fmt.Sprintf("%d", hash.Sum32())
	toMatch := c.GetHeader("If-None-Match")
	if etag == toMatch {
		c.Status(http.StatusNotModified)
		return
	}

	c.Header("ETag", etag)
	c.Header("Cache-Control", "no-cache")
	c.IndentedJSON(http.StatusOK, tasks)
}

// Return specific task.
// Cache using crc32 ETag.
func GetTask(c *gin.Context) {
	conn, err := db.GetConnection(c)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer conn.Close()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ClientError(c, err)
		return
	}

	rows, err := conn.QueryContext(c, "SELECT id, title, done, created_ts FROM tasks WHERE id=?;", id)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer rows.Close()

	// Fetch and calculate ETag
	task := models.Task{}
	if !rows.Next() {
		c.String(http.StatusNotFound, "Not Found")
		return
	}
	err = rows.Scan(&task.ID, &task.Title, &task.Done, &task.CreatedTS)
	if err != nil {
		InternalServerError(c, err)
		return
	}

	// Hash
	hash := crc32.NewIEEE()
	task.HashWrite32(&hash)

	// Compare client and server ETags
	etag := fmt.Sprintf("W/%d", hash.Sum32())
	toMatch := c.GetHeader("If-None-Match")
	if etag == toMatch {
		c.Status(http.StatusNotModified)
		return
	}

	c.Header("ETag", etag)
	c.Header("Cache-Control", "no-cache")
	c.IndentedJSON(http.StatusOK, task)
}

// Create new task.
// Returns 201 Created
func PostTask(c *gin.Context) {
	// Parse JSON body
	var newTask models.Task
	if err := c.BindJSON(&newTask); err != nil {
		InternalServerError(c, err)
		return
	}
	newTask.ID = uuid.New()
	newTask.CreatedTS = time.Now().UTC().Round(1 * time.Second)

	// Get SQL connection
	conn, err := db.GetConnection(c)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer conn.Close()

	// Insert task
	sql := "INSERT INTO tasks " +
		"(id, title, done, created_ts)" +
		"VALUES (?, ?, ?, ?)"
	stmt, err := conn.PrepareContext(c, sql)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(newTask.ID, newTask.Title, newTask.Done, newTask.CreatedTS)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	_ = res

	// Return 201 Created
	c.Header("Location", fmt.Sprintf("/tasks/%s", newTask.ID.String()))
	c.Status(http.StatusCreated)
}

func UpdateTask(c *gin.Context) {

}

func DeleteTask(c *gin.Context) {
	conn, err := db.GetConnection(c)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	defer conn.Close()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		ClientError(c, err)
		return
	}

	res, err := conn.ExecContext(c, "DELETE FROM tasks WHERE id=?;", id)
	if err != nil {
		InternalServerError(c, err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		InternalServerError(c, err)
		return
	}
	if rows < 1 {
		c.String(http.StatusNotFound, "Not Found")
		return
	}

	c.Header("Cache-Control", "no-cache")
	c.Status(http.StatusNoContent)
}
