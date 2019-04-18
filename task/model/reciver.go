package model

import (
	"database/sql"
)

type ReceiveInfo struct {
	taskID int
	userID int
}

const (
	mysqlReceiveCreateTable = iota

	mysqlReceiveInsertInfo
	mysqlReceiveDeleteInfo

	mysqlReceiveQueryUser
	mysqlReceiveQueryTask
)

var (
	receiveSQLString = []string{
		`CREATE TABLE IF NOT EXISTS receiveInfo (
			taskId INT UNSIGNED NOT NULL , 
			userId INT UNSIGNED NOT NULL 
		  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,

		`INSERT INTO receiveInfo (userId,taskId) VALUES (?,?)`,
		`DELETE FROM receiveInfo WHERE taskId = ? , userId = ?`,

		`SELECT userId FROM receiveInfo WHERE taskId = ? LOCK IN SHARE MODE`,
		`SELECT taskId FROM receiveInfo WHERE userId = ? LOCK IN SHARE MODE`,
	}
)

func CreateReceiveTable(db *sql.DB) error {
	_, err := db.Exec(receiveSQLString[mysqlReceiveCreateTable])

	if err != nil {
		return err
	}

	return nil

}

func InsertReceiverInfo(db *sql.DB, userID, taskID int) error {
	result, err := db.Exec(receiveSQLString[mysqlReceiveInsertInfo], userID, taskID)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}

func DeleteReceiveInfo(db *sql.DB, userID, taskID int) error {
	_, err := db.Exec(receiveSQLString[mysqlReceiveDeleteInfo], userID, taskID)

	if err != nil {
		return err
	}

	return nil
}

func QueryUserIDByTaskID(db *sql.DB, taskID int) (*[]int, error) {
	rows, err := db.Query(receiveSQLString[mysqlReceiveQueryUser], taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersID []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		usersID = append(usersID, userID)
	}

	return &usersID, nil
}

func QueryTaskByUserID(db *sql.DB, userID int) (*[]int, error) {
	rows, err := db.Query(receiveSQLString[mysqlReceiveQueryTask], userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskesID []int
	for rows.Next() {
		var taskID int
		if err := rows.Scan(&taskID); err != nil {
			return nil, err
		}

		taskesID = append(taskesID, taskID)
	}

	return &taskesID, nil

}
