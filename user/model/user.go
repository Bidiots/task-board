package model

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

const (
	mysqlUserCreateTable = iota
	mysqlUserInsert
	mysqlUserDeleteByID
	mysqlUserInfoByID
	mysqlUserInfoByName
)

var (
	errInvalidInsert = errors.New("insert user:insert affected 0 rows")

	UserSQLString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			userId    INT NOT NULL AUTO_INCREMENT,
			name        VARCHAR(100) UNIQUE DEFAULT NULL,
			password	VARCHAR(40)  DEFAULT NULL,
			PRIMARY KEY (userId)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8_mb4`,
		`INSERT INTO  %s (name,password) VALUES (?,?)`,
		`DELETE FROM %s WHERE userId = ? LIMIT 1`,
		`SELECT * FROM %s WHERE userId = ? LIMIT 1 LOCK IN SHARE MODE`,
		`SELECT password FROM %s WHERE name=? LIMIT 1 LOCK IN SHARE MODE`,
	}
)

func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(UserSQLString[mysqlUserCreateTable], tableName)
	_, err := db.Exec(sql)

	return err
}

func InsertUser(db *sql.DB, tableName string, name string, password string) (int, error) {
	sql := fmt.Sprintf(UserSQLString[mysqlUserInsert], tableName)
	result, err := db.Exec(sql, name, password)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	UserID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(UserID), nil
}

func DeleteByID(db *sql.DB, tableName string, id int) error {
	sql := fmt.Sprintf(UserSQLString[mysqlUserDeleteByID], tableName)
	_, err := db.Exec(sql, id)
	if err != nil {
		return err
	}
	return nil
}

func InfoByID(db *sql.DB, tableName string, id int) (*User, error) {
	var user User

	sql := fmt.Sprintf(UserSQLString[mysqlUserInfoByID], tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Password); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func InfoPasswordByName(db *sql.DB, tableName string, name string) (string, error) {
	sql := fmt.Sprintf(UserSQLString[mysqlUserInfoByName], tableName)
	rows, err := db.Query(sql, name)
	if err != nil {
		return "", err
	}

	var password string
	for rows.Next() {
		err = rows.Scan(&password)
		if err != nil {
			return "", err
		}

	}

	return password, nil
}
