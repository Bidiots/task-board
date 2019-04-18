package model

import (
	"database/sql"
	"time"
)

const (
	mysqlPermissionCreateTable = iota
	mysqlPermissionInstert
	mysqlPermissionDelete
	mysqlPermissonGetRole
	mysqlPermissonGetAll
	mysqlPermissonGetMap
)

type permission struct {
	URL       string
	RoleID    uint32
	CreatedAt time.Time
}

var (
	permissionSQLString = []string{
		`CREATE TABLE IF NOT EXISTS permission (
			url			VARCHAR(255) NOT NULL DEFAULT ' ',
			roleId		MEDIUMINT UNSIGNED NOT NULL,
			createdAt 	DATETIME UNIQUE DEFAULT NULL,
			PRIMARY KEY (url,roleId)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=1000`,
		`INSERT INTO permission(url,roleId) VALUES (?,?)`,
		`DELETE FROM permission WHERE roleId = ? AND url = ? LIMIT 1`,
		`SELECT * FROM permission, role WHERE url = ? LOCK IN SHARE MODE`,
		`SELECT * FROM permission LOCK IN SHARE MODE`,
		`SELECT permission.roleId FROM permission, role WHERE permission.url = ? AND permission.roleId = role.id LOCK IN SHARE MODE`,
	}
)

//CreatePermissionTable -
func CreatePermissionTable(db *sql.DB) error {
	_, err := db.Exec(permissionSQLString[mysqlPermissionCreateTable])
	return err
}

//InsertURLPermission -
func InsertURLPermission(db *sql.DB, rid int, url string) error {
	_, err := db.Exec(permissionSQLString[mysqlPermissionInstert], url, rid)
	return err
}

//DeleteURLPermission -
func DeleteURLPermission(db *sql.DB, rid int, url string) error {
	_, err := db.Exec(permissionSQLString[mysqlPermissionDelete], rid, url)
	return err
}

//InfoURLPermissions -
func InfoURLPermissions(db *sql.DB, url string) (*permission, error) {
	var (
		roleID    uint32
		createdAt time.Time
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetAll])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&url, &roleID, &createdAt); err != nil {
			return nil, err
		}
		result := &permission{
			URL:       url,
			RoleID:    roleID,
			CreatedAt: createdAt,
		}
		return result, nil
	}
	return nil, nil
}

//InfoPermissions -
func InfoPermissions(db *sql.DB) ([]*permission, error) {
	var (
		roleID    uint32
		url       string
		createdAt time.Time

		result []*permission
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetAll])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&url, &roleID, &createdAt); err != nil {
			return nil, err
		}
		data := &permission{
			URL:       url,
			RoleID:    roleID,
			CreatedAt: createdAt,
		}
		result = append(result, data)
	}
	return result, nil
}

//URLPermissionsMap -
func URLPermissionsMap(db *sql.DB, url string) (map[int]bool, error) {
	var (
		roleID int
		result = make(map[int]bool)
	)

	rows, err := db.Query(permissionSQLString[mysqlPermissonGetRole], url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&roleID); err != nil {
			return nil, err
		}
		result[roleID] = true
	}
	return result, nil
}
