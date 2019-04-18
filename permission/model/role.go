package model

import (
	"database/sql"
	"errors"
	"time"
)

type role struct {
	ID       uint32
	Name     string
	Intro    string
	CreateAt time.Time
}

const (
	mysqlRoleCreateTable = iota
	mysqlRoleInsert
	mysqlRoleModify
	mysqlRoleGetList
	mysqlRoleGetByID
)

var (
	errInvalidMysql = errors.New("affected 0 rows")
	//roleSQLString -
	roleSQLString = []string{
		`CREATE TABLE IF NOT EXISTS role (
			id 	        INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name		VARCHAR(255) UNIQUE NOT NULL DEFAULT ' ',
			intro		VARCHAR(255) NOT NULL DEFAULT ' ',
			createdAt 	DATETIME UNIQUE DEFAULT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1000`,
		`INSERT INTO role(name,intro,createAt) VALUES (?,?,?)`,
		`UPDATE role SET name = ?,intro = ? WHERE id = ? LIMIT 1`,
		`SELECT * FROM role LOCK IN SHARE MODE`,
		`SELECT * FROM role WHERE id = ? LOCK IN SHARE MODE`,
	}
)

//CreateRoleTable -
func CreateRoleTable(db *sql.DB) error {
	_, err := db.Exec(roleSQLString[mysqlRoleCreateTable])

	if err != nil {
		return err
	}

	return nil
}

//CreateRole -
func CreateRole(db *sql.DB, name, intro string) error {
	result, err := db.Exec(roleSQLString[mysqlRoleInsert], name, intro, time.Now())
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

//ModifyRole -
func ModifyRole(db *sql.DB, id int, name, intro string) error {
	_, err := db.Exec(roleSQLString[mysqlRoleModify], name, intro, id)

	if err != nil {
		return err
	}

	return nil
}

//InfoRoleList -
func InfoRoleList(db *sql.DB) ([]*role, error) {
	var (
		id       uint32
		name     string
		intro    string
		createAt time.Time
		roles    []*role
	)
	rows, err := db.Query(roleSQLString[mysqlRoleGetList])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &intro, &createAt); err != nil {
			return nil, err
		}

		r := role{
			ID:       id,
			Name:     name,
			Intro:    intro,
			CreateAt: createAt,
		}

		roles = append(roles, &r)
	}

	return roles, nil
}

//GetRoleByID -
func GetRoleByID(db *sql.DB, id int) (*role, error) {
	var (
		r role
	)

	err := db.QueryRow(roleSQLString[mysqlRoleGetByID], id).Scan(&r.ID, &r.Name, &r.Intro, &r.CreateAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
