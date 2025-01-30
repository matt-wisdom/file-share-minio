package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type dbController struct {
	db *sql.DB
}

// Add methods to dbController
func (dbC *dbController) init_db() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	dbC.db = db

}

func (dbC *dbController) createUserDB(username, email string) int {
	var id int
	query := "INSERT into users (email, username) VALUES ($1, $2) RETURNING user_id"
	err := dbC.db.QueryRow(query, email, username).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}

func (dbC *dbController) deleteUserDB(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := dbC.db.Exec(query, userID)
	return err
}

func (dbC *dbController) createFileDB(filename, bucketId string, ownerID int) (int, error) {
	var id int
	query := "INSERT into files (file_name, object_name, owner_id) VALUES ($1, $2, $3) RETURNING file_id"
	err := dbC.db.QueryRow(query, filename, bucketId, ownerID).Scan(&id)
	return id, err
}

func (dbC *dbController) deleteFileDB(fileID int) error {
	query := "DELETE FROM files WHERE id = $1"
	_, err := dbC.db.Exec(query, fileID)
	// Delete shares
	query = "DELETE FROM shares WHERE file_id = $1"
	_, _ = dbC.db.Exec(query, fileID)
	return err
}

func (dbC *dbController) shareFileDB(fileID, userID int) error {
	query := "INSERT into shares (file_id, shared_with) VALUES ($1, $2)"
	_, err := dbC.db.Exec(query, fileID, userID)
	return err
}

func (dbC *dbController) unshareFileDB(fileID, userID int) error {
	query := "DELETE FROM shares WHERE file_id = $1 AND shared_with = $2"
	_, err := dbC.db.Exec(query, fileID, userID)
	return err
}

func (dbC *dbController) getUserFilesDB(userID int) ([]file, error) {
	var files []file
	query := "SELECT file_name FROM files WHERE owner_id = $1"
	rows, err := dbC.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f file
		err := rows.Scan(&f.file_id, &f.object_name, &f.file_name, &f.owner_id, &f.created_at)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (dbC *dbController) getSharedFilesDB(userID int) ([]file, error) {
	var files []file
	query := "SELECT f.file_id, f.object_name, f.file_name, f.owner_id, f.created_at FROM files f JOIN shares s ON f.file_id = s.file_id WHERE s.shared_with = $1"
	rows, err := dbC.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f file
		err := rows.Scan(&f.file_id, &f.object_name, &f.file_name, &f.owner_id, &f.created_at)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (dbC *dbController) getFileSharedFromUserDB(receiverID, sharerID int) ([]file, error) {
	var files []file
	query := "SELECT f.file_id, f.object_name, f.file_name, f.owner_id, f.created_at FROM files f JOIN shares s ON f.file_id = s.file_id WHERE s.shared_with = $1 AND f.owner_id = $2"
	rows, err := dbC.db.Query(query, receiverID, sharerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f file
		err := rows.Scan(&f.file_id, &f.object_name, &f.file_name, &f.owner_id, &f.created_at)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}
