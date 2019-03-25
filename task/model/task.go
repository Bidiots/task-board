package model

import (
	"TEST/user/model"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID          int        `json:"taskId`
	Name        string     `json:name`
	Description string     `json:description`
	CreateTime  time.Time  `json:createTime`
	Reciver     model.User `json:receiver`
	Poster      model.User `json:poster`
}

const (
	mysqlTaskCreateTable = iota
	mysqlTaskInsert
	mysqlTaskInfoByID
	mysqlTaskDeleteByID
	mysqlTaskInfoAll
	mysqlTaskInfoDescripty
	mysqlTaskPosterByID
	mysqlTaskUpdateByID
)

var (
	errInvalidInsert = errors.New("insert task:insert affected 0 rows")
	TaskSQLString    = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			taskId    INT NOT NULL AUTO_INCREMENT,
			name        VARCHAR(100) UNIQUE DEFAULT NULL,
			description VARCHAR(255) UNIQUE DEFAULT NULL,
			createTime  DATATIME UNIQUE DEFAULT NULL,
			receiver	VARCHAR(100)	
			poster		VARCHAR(100)	UNIQUE DEFAULT NULL,
			PRIMARY KEY (TaskId)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8`,
		`INSERT INTO  %s (name,description,createtime,poster) VALUES (?,?)`,
		`SELECT * FROM %s WHERE taskId = ? LIMIT 1 LOCK IN SHARE MODE`,
		`DELETE FROM %s WHERE taskId = ? LIMIT 1`,
		`SELETE * FROM %s`,
		`SELETE description FROM %s WHERE taskId = ? LIMIT 1`,
		`SELETE poster FROM %s WHERE taskId = ? LIMIT 1`,
		`UPDATE %s SET description=? WHERE taskId = ? LIMIT 1`,
	}
)

func UpdateDescriptionByID(db *sql.DB, tableName string, id int, descripty string) error {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskUpdateByID], tableName)
	_, err := db.Exec(sql, descripty, id)
	if err != nil {
		return err
	}
	return nil
}
func InfoPosterNameByID(db *sql.DB, tableName string, id int) (string, error) {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskPosterByID], tableName)
	rows, err := db.Query(sql, id)
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
func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskCreateTable], tableName)
	_, err := db.Exec(sql)

	return err
}
func InsertTask(db *sql.DB, tableName string, name string, description string, createTime time.Time, user string) (int, error) {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskInsert], tableName)
	result, err := db.Exec(sql, name, description, createTime)
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
func InfoByID(db *sql.DB, tableName string, id int) (*Task, error) {
	var task Task
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskInfoByID], tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.CreateTime, &task.Reciver, &task.Poster); err != nil {
			return nil, err
		}
	}
	return &task, nil
}
func DeleteByID(db *sql.DB, tableName string, id int) error {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskDeleteByID], tableName)
	_, err := db.Exec(sql, id)
	return err
}
func InfoAllTask(db *sql.DB, tableName string) ([]Task, error) {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskInfoAll], tableName)
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	i := 0
	task := make([]Task, 100, 100)
	for rows.Next() {
		i++
		if err := rows.Scan(&task[i].ID, &task[i].Name, &task[i].Description, &task[i].CreateTime, &task[i].Reciver, &task[i].Poster); err != nil {
			return nil, err
		}
	}
	return task, err
}
func InfoDescription(db *sql.DB, tableName string, id int) (string, error) {
	sql := fmt.Sprintf(TaskSQLString[mysqlTaskInfoDescripty], tableName)
	rows, err := db.Query(sql, id)
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
