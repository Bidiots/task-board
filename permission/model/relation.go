package model

import (
	"database/sql"
	"errors"
	"time"
)

type relation struct {
	AdminID uint32
	RoleID  uint32
}

const (
	mysqlRelationCreateTable = iota
	mysqlRelationInsert
	mysqlRelationDelete
	mysqlRelationRoleMap
)

var (
	errAdminInactive = errors.New("the admin is not activated")
	errRoleInactive  = errors.New("the role is not activated")

	relationSqlString = []string{
		`CREATE TABLE IF NOT EXISTS relation (
			adminId 	BIGINT UNSIGNED NOT NULL,
			roleId	INT UNSIGNED NOT NULL,
			createdAt 	DATETIME UNIQUE DEFAULT NULL ,
			PRIMARY KEY (adminId,roleId)
		) ENGINE=InnoDB  DEFAULT CHARSET=utf8`,
		`INSERT INTO relation(adminId,roleId,createdAt) VALUES (?,?,?)`,
		`DELETE FROM relation WHERE adminId = ? AND roleId = ? LIMIT 1`,
		`SELECT relation.roleId FROM relation, role WHERE relation.adminId = ? AND relation.roleId = role.id LOCK IN SHARE MODE`,
	}
)

func GetRoleMap(db *sql.DB, userID int) (map[int]bool, error) {
	var (
		roleID int
		result = make(map[int]bool)
	)
	rows, err := db.Query(roleSqlString[mysqlRelationRoleMap], userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&roleID); err != nil {
			return nil, err
		}
		result[roleID] = true
	}
	return result, nil
}

func CreateRelationTable(db *sql.DB) error {
	_, err := db.Exec(relationSqlString[mysqlRelationCreateTable])
	return err
}

func InsertRelation(db *sql.DB, aid, rid int) error {

	_, err := db.Exec(relationSqlString[mysqlRelationInsert], aid, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func DeleteRelation(db *sql.DB, aid, rid int) error {
	_, err := db.Exec(relationSqlString[mysqlRelationDelete], aid, rid)

	return err
}
