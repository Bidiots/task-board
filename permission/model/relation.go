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
			admin_id 	BIGINT UNSIGNED NOT NULL,
			role_id		INT UNSIGNED NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (admin_id,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO relation(admin_id,role_id,created_at) VALUES (?,?,?,?)`,
		`DELETE FROM relation WHERE admin_id = ? AND role_id = ? LIMIT 1`,
		`SELECT relation.role_id FROM relation, role WHERE relation.admin_id = ? AND role.active = true AND relation.role_id = role.id LOCK IN SHARE MODE`,
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

	_, err := db.Exec(relationSqlString[mysqlRelationInsert], aid, rid, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func DeleteRelation(db *sql.DB, aid, rid int) error {
	_, err := db.Exec(relationSqlString[mysqlRelationDelete], aid, rid)

	return err
}
