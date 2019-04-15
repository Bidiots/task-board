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

	roleSqlString = []string{
		`CREATE TABLE IF NOT EXISTS role (
			id 	        INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name		VARCHAR(255) UNIQUE NOT NULL DEFAULT ' ',
			intro		VARCHAR(255) NOT NULL DEFAULT ' ',
			createdAt 	DATETIME UNIQUE DEFAULT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8`,
		`INSERT INTO role(name,intro,createAt) VALUES (?,?,?)`,
		`UPDATE role SET name = ?,intro = ? WHERE id = ? LIMIT 1`,
		`SELECT * FROM role LOCK IN SHARE MODE`,
		`SELECT * FROM role WHERE id = ? LOCK IN SHARE MODE`,
	}
)

func CreateRoleTable(db *sql.DB) error {
	_, err := db.Exec(roleSqlString[mysqlRoleCreateTable])

	if err != nil {
		return err
	}

	return nil
}

func CreateRole(db *sql.DB, name, intro string) error {
	result, err := db.Exec(roleSqlString[mysqlRoleInsert], name, intro, time.Now())
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

func ModifyRole(db *sql.DB, id int, name, intro string) error {
	_, err := db.Exec(roleSqlString[mysqlRoleModify], name, intro, id)

	if err != nil {
		return err
	}

	return nil
}

func InfoRoleList(db *sql.DB) (*[]role, error) {
	var (
		id       uint32
		name     string
		intro    string
		createAt time.Time
		roles    []role
	)
	rows, err := db.Query(roleSqlString[mysqlRoleGetList])
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

		roles = append(roles, r)
	}

	return &roles, nil
}

func GetRoleByID(db *sql.DB, id int) (*role, error) {
	var (
		r role
	)

	err := db.QueryRow(roleSqlString[mysqlRoleGetByID], id).Scan(&r.ID, &r.Name, &r.Intro, &r.CreateAt)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
