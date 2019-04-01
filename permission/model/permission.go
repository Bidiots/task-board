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
	Url       string
	RoleID    uint32
	CreatedAt time.Time
}

var (
	permissionSqlString = []string{
		`CREATE TABLE IF NOT EXISTS permission (
			url			VARCHAR(512) NOT NULL DEFAULT ' ',
			role_id		MEDIUMINT UNSIGNED NOT NULL,
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (url,role_id)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO permission(url,role_id) VALUES (?,?)`,
		`DELETE FROM permission WHERE role_id = ? AND url = ? LIMIT 1`,
		`SELECT * FROM permission, role WHERE url = ? `,
		`SELECT * FROM permission LOCK IN SHARE MODE`,
		`SELECT permission.role_id FROM permission, role WHERE permission.url = ? AND permission.role_id = role.id LOCK IN SHARE MODE`,
	}
)

func CreatePermissionTable(db *sql.DB) error {
	_, err := db.Exec(permissionSqlString[mysqlPermissionCreateTable])
	return err
}

func InsertURLPermission(db *sql.DB, rid int, url string) error {
	_, err := db.Exec(permissionSqlString[mysqlPermissionInstert], url, rid)
	return err
}

func DeleteURLPermission(db *sql.DB, rid int, url string) error {
	_, err := db.Exec(permissionSqlString[mysqlPermissionDelete], rid, url)
	return err
}

func InfoURLPermissions(db *sql.DB, url string) (*permission, error) {
	var (
		roleID    uint32
		createdAt time.Time
	)

	rows, err := db.Query(permissionSqlString[mysqlPermissonGetAll])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&url, &roleID, &createdAt); err != nil {
			return nil, err
		}
		result := &permission{
			Url:       url,
			RoleID:    roleID,
			CreatedAt: createdAt,
		}
		return result, nil
	}
	return nil, nil
}

func InfoPermissions(db *sql.DB) (*[]permission, error) {
	var (
		roleID    uint32
		url       string
		createdAt time.Time

		result []permission
	)

	rows, err := db.Query(permissionSqlString[mysqlPermissonGetAll])
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
func URLPermissionsMap(db *sql.DB, url string) (map[int]bool, error) {
	var (
		roleID int
		result = make(map[int]bool)
	)

	rows, err := db.Query(permissionSqlString[mysqlPermissonGetRole], url)
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
