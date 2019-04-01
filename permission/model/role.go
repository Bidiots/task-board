package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ServiceProvider struct{}

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
	mysqlRoleModifyActive
	mysqlRoleGetList
	mysqlRoleGetByID
)

var errInvalidMysql = errors.New("affected 0 rows")

var roleSqlString = []string{
	`CREATE TABLE IF NOT EXISTS %s (
			id 	        INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name		VARCHAR(512) UNIQUE NOT NULL DEFAULT ' ',
			intro		VARCHAR(512) NOT NULL DEFAULT ' ',
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	`INSERT INTO %s(name,intro,active) VALUES (?,?,?)`,
	`UPDATE %s SET name = ?,intro = ? WHERE id = ? LIMIT 1`,
	`UPDATE %s SET active = ? WHERE id = ? LIMIT 1`,
	`SELECT * FROM %s LOCK IN SHARE MODE`,
	`SELECT * FROM %s WHERE id = ? AND active = true LOCK IN SHARE MODE`,
}

func CreateRoleTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(roleSqlString[mysqlPermissionCreateTable], tableName)
	_, err := db.Exec(roleSqlString[mysqlRoleCreateTable])
	return err
}

func CreateRole(db *sql.DB, tableName, name, intro string) error {
	sql := fmt.Sprintf(roleSqlString[mysqlRoleInsert], tableName)
	result, err := db.Exec(sql, name, intro, true)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidMysql
	}

	return nil
}

func ModifyRole(db *sql.DB,tableName string ,id uint32, ,name, intro string) error {
	sql:=fmt.Sprintf(roleSqlString[mysqlRoleModify],tableName)
	_, err := db.Exec(sql, name, intro, id)

	return err
}

// ModifyRoleActive modify role active.
func ModifyRoleActive(db *sql.DB,tableName string, id uint32, active bool) error {
	sql:=fmt.Sprintf(roleSqlString[mysqlRoleModifyActive])
	_, err := db.Exec(roleSqlString[mysqlRoleModifyActive], active, id)

	return err
}

// RoleList get all role information.
func RoleList(db *sql.DB,tableName string) (*[]role, error) {
	
	var (
		id       uint32
		name     string
		intro    string
		active   bool
		createAt time.Time

		roles    []role
	)
	sql:=fmt.Sprintf(roleSqlString[mysqlRoleGetList],tableName)
	rows, err := db.Query(roleSqlString[mysqlRoleGetList])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &intro, &active, &createAt); err != nil {
			return nil, err
		}

		r := role{
			ID:       id,
			Name:     name,
			Intro:    intro,
			Active:   active,
			CreateAt: createAt,
		}

		roles = append(roles, r)
	}

	return &roles, nil
}

// GetRoleByID get role by id.
func GetRoleByID(db *sql.DB,tableName string, id uint32) (*role, error) {
	sql:=fmt.Sprintf(roleSqlString[mysqlRoleGetByID],tableName)
	var (
		r role
	)
	err := db.QueryRow(sql, id).Scan(&r.ID, &r.Name, &r.Intro, &r.CreateAt)
	return &r, err
}
