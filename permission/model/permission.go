package model

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	mysqlPermissionCreateTable = iota
	mysqlPermissionInstert
	mysqlPermissionDelete
	mysqlPermissonGetRole
	mysqlPermissonGetAll
)

type permission struct {
	Url       string
	RoleID    uint32
	CreatedAt time.Time
}

var (
	permissionSqlString = []string{
		`CREATE TABLE IF NOT EXISTS %s (
			url			VARCHAR(512) NOT NULL DEFAULT ' ',
			role_id		MEDIUMINT UNSIGNED NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (url,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO %s(url,role_id) VALUES (?,?)`,
		`DELETE FROM %s WHERE role_id = ? AND url = ? LIMIT 1`,
		`SELECT permission.role_id FROM %s, admin.role WHERE permission.url = ? AND role.active = true AND permission.role_id = role.id LOCK IN SHARE MODE`,
		`SELECT * FROM %s LOCK IN SHARE MODE`,
	}
)

// CreatePermissionTable create permission table.
func CreatePermissionTable(db *sql.DB, tableName string) error {
	sql := fmt.Sprintf(permissionSqlString[mysqlPermissionCreateTable], tableName)
	_, err := db.Exec(sql)
	return err
}

// AddPermission create an associated record of the specified URL and role.
func AddURLPermission(db *sql.DB, tableName string, rid uint32, url string) error {
	sql := fmt.Sprintf(permissionSqlString[mysqlPermissionInstert], tableName)
	role, err := sp.GetRoleByID(db, rid)
	if err != nil {
		return err
	}

	_, err = db.Exec(sql, url, rid)
	if err != nil {
		return err
	}

	return nil
}

// RemovePermission remove the associated records of the specified URL and role.
func RemoveURLPermission(db *sql.DB, tableName string, rid uint32, url string) error {
	sql := fmt.Sprintf(permissionSqlString[mysqlPermissionDelete], tableName)
	role, err := sp.GetRoleByID(db, rid)
	if err != nil {
		return err
	}

	_, err = db.Exec(sql, rid, url)
	if err != nil {
		return err
	}

	return nil
}

// URLPermissions lists all the roles of the specified URL.
func URLPermissions(db *sql.DB, tableName string, url string) (map[uint32]bool, error) {
	sql := fmt.Sprintf(permissionSqlString[mysqlPermissonGetRole], tableName)
	var (
		roleID uint32
		result = make(map[uint32]bool)
	)

	rows, err := db.Query(sql, url)
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

// Permissions lists all the roles.
func Permissions(db *sql.DB, tableName string) (*[]permission, error) {
	sql := fmt.Sprintf(permissionSqlString[mysqlPermissonGetAll], tableName)
	var (
		roleID    uint32
		url       string
		createdAt time.Time

		result []permission
	)

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&url, &roleID, &createdAt); err != nil {
			return nil, err
		}
		data := permission{
			Url:       url,
			RoleID:    roleID,
			CreatedAt: createdAt,
		}
		result = append(result, data)
	}
	return &result, nil
}
