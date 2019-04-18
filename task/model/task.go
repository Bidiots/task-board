package model

import (
	"database/sql"
	"errors"
	"time"
)

//Task -
type Task struct {
	ID          int       `json:"taskId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"createTime"`
	Poster      string    `json:"poster"`
}

const (
	mysqlTaskCreateTable = iota

	mysqlTaskInsert
	mysqlTaskDeleteByID
	mysqlTaskUpdateByID

	mysqlTaskInfoAll
	mysqlTaskInfoByID
	mysqlTaskInfoDescripty
	mysqlTaskInfoPosterByID
	mysqlTaskInfoReceiverById
)

var (
	errInvalidInsert = errors.New("insert task:insert affected 0 rows")
	//TaskSQLString -
	TaskSQLString = []string{
		`CREATE TABLE IF NOT EXISTS tasks (
			taskId INT UNSIGNED NOT NULL AUTO_INCREMENT, 
			name VARCHAR(100) UNIQUE DEFAULT "" NOT NULL, 
			description VARCHAR(255) UNIQUE DEFAULT "" NOT NULL,
			createTime DATETIME UNIQUE DEFAULT NULL,
			poster VARCHAR(100) UNIQUE DEFAULT NULL, 
			PRIMARY KEY (taskId)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1000`,

		`INSERT INTO tasks (name,description,createtime,poster) VALUES (?,?,?,?)`,
		`DELETE FROM tasks WHERE taskId = ? LIMIT 1`,
		`UPDATE tasks SET description=? WHERE taskId = ? LIMIT 1`,

		`SELECT * FROM tasks LOCK IN SHARE MODE`,
		`SELECT * FROM tasks WHERE taskId = ? LIMIT 1 LOCK IN SHARE MODE`,
		`SELECT description FROM tasks WHERE taskId = ? LIMIT 1 LOCK IN SHARE MODE`,
		`SELECT poster FROM tasks WHERE taskId = ? LIMIT 1 LOCK IN SHARE MODE`,
		`SELECT tasks FROM tasks WHERE taskId = ? LIMIT 1 LOCK IN SHARE MODE`,
	}
)

//CreateTaskTable -
func CreateTaskTable(db *sql.DB) error {
	_, err := db.Exec(TaskSQLString[mysqlTaskCreateTable])
	if err != nil {
		return err
	}

	return nil
}

// InsertTask -
func InsertTask(db *sql.DB, name string, description string, createTime time.Time, user string) (int, error) {
	result, err := db.Exec(TaskSQLString[mysqlTaskInsert], name, description, createTime, user)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	TaskID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(TaskID), err

}

//DeleteByID -
func DeleteByID(db *sql.DB, id int) error {
	_, err := db.Exec(TaskSQLString[mysqlTaskDeleteByID], id)

	if err != nil {
		return err
	}

	return nil
}

//InfoPosterNameByID -
func InfoPosterNameByID(db *sql.DB, id int) (string, error) {
	rows, err := db.Query(TaskSQLString[mysqlTaskInfoPosterByID], id)
	if err != nil {
		return "", err
	}

	var msg string
	for rows.Next() {
		if err = rows.Scan(&msg); err != nil {
			return "", err
		}
	}

	return msg, nil
}

//UpdateDescriptionByID -
func UpdateDescriptionByID(db *sql.DB, id int, descripty string) error {
	_, err := db.Exec(TaskSQLString[mysqlTaskUpdateByID], descripty, id)

	if err != nil {
		return err
	}

	return nil
}

//InfoByID -
func InfoByID(db *sql.DB, id int) (*Task, error) {
	var task Task
	rows, err := db.Query(TaskSQLString[mysqlTaskInfoByID], id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.CreateTime, &task.Poster); err != nil {
			return nil, err
		}
	}

	return &task, nil
}

//InfoAllTask -
func InfoAllTask(db *sql.DB) (tasks []Task, err error) {
	var (
		id          int
		name        string
		description string
		createTime  time.Time
		poster      string
	)

	rows, err := db.Query(TaskSQLString[mysqlTaskInfoAll])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &description, &createTime, &poster); err != nil {
			return nil, err
		}

		task := Task{
			ID:          id,
			Name:        name,
			Description: description,
			CreateTime:  createTime,
			Poster:      poster,
		}

		tasks = append(tasks, task)
	}

	return tasks, err
}

//InfoDescription -
func InfoDescription(db *sql.DB, id int) (string, error) {
	rows, err := db.Query(TaskSQLString[mysqlTaskInfoDescripty], id)
	if err != nil {
		return "", nil
	}

	var msg string
	for rows.Next() {
		if err := rows.Scan(&msg); err != nil {
			return "", err
		}

	}

	return msg, err

}
