package model

import (
	"database/sql"
	"errors"
	"fmt"
)

type Admin struct {
	ID       int64  `json:"id"`
	Name     string `json:"adminname"`
	Password string `json:"password"`
}

const (
	mysqlAdminCreateTable = iota
	mysqlAdminInsert
	mysqlAdminInfoByID
	mysqlAdminDeleteByID
	mysqlAdminInfoByName
)

var (
	errInvalidInsert = errors.New("insert admin:insert affected 0 rows")

	AdminSQLString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			AdminId    INT NOT NULL AUTO_INCREMENT,
			name        VARCHAR(100) UNIQUE DEFAULT NULL,
			password	VARCHAR(40) UNIQUE DEFAULT NULL,
			PRIMARY KEY (AdminId)
		)ENGINE=InnoDB DEFAULT CHARSET=utf8`,
		`INSERT INTO  %s (name,password) VALUES (?,?)`,
		`SELECT * FROM %s WHERE adminid = ? LIMIT 1 LOCK IN SHARE MODE`,
		`DELETE FROM %s WHERE adminid = ? LIMIT 1`,
		`SELECT password FROM ? WHERE name=? LIMIT 1`,
	}
)

func CreateTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(AdminSQLString[mysqlAdminCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}
func InsertAdmin(db *sql.DB, tableName string, name string, password string) (int, error) {
	sql := fmt.Sprintf(AdminSQLString[mysqlAdminInsert], tableName)
	result, err := db.Exec(sql, name, password)
	if err != nil {
		return 0, err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return 0, errInvalidInsert
	}

	AdminID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(AdminID), nil
}
func InfoByID(db *sql.DB, tableName string, id int) (*Admin, error) {
	var Admin Admin

	sql := fmt.Sprintf(AdminSQLString[mysqlAdminInfoByID], tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&Admin.ID, &Admin.Name, &Admin.Password); err != nil {
			return nil, err
		}
	}

	return &Admin, nil
}
func DeleteByID(db *sql.DB, tableName string, id int) error {
	sql := fmt.Sprintf(AdminSQLString[mysqlAdminDeleteByID], tableName)
	_, err := db.Exec(sql, id)
	return err
}
func InfoPasswordByName(db *sql.DB, tableName string, name string) (string, error) {
	sql := fmt.Sprintf(AdminSQLString[mysqlAdminInfoByName], tableName, name)
	rows, err := db.Query(sql)
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
	return password, err
}
