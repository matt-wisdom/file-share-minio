package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type dbController struct {
	db *sql.DB
}

// Add methods to dbController
func (dbC *dbController) init_db() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	dbC.db = db

}

func (dbC *dbController) close_db() {
	dbC.db.Close()
}

func (dbC *dbController) getUserDB(userID int) (userModel, error) {
	var user userModel
	query := "SELECT user_id, username, email FROM users WHERE user_id = $1"
	err := dbC.db.QueryRow(query, userID).Scan(&user.user_id, &user.username, &user.email)
	log.Print("GetUser:", err, "\n")
	if err != nil {
		return user, err
	}
	return user, nil
}

func (dbC *dbController) getUserByEmailDB(email string) (userModel, error) {
	var user userModel
	query := "SELECT user_id, username, email FROM users WHERE email = $1"
	err := dbC.db.QueryRow(query, email).Scan(&user.user_id, &user.username, &user.email)
	log.Print("GetUserByEmail:", err, "\n")
	if err != nil {
		return user, err
	}
	return user, nil
}

func (dbC *dbController) getUserByUsernameDB(username string) (userModel, error) {
	var user userModel
	query := "SELECT user_id, username, email FROM users WHERE username = $1"
	err := dbC.db.QueryRow(query, username).Scan(&user.user_id, &user.username, &user.email)
	log.Print("GetUsername:", err, "\n")
	if err != nil {
		return user, err
	}
	return user, nil
}

func (dbC *dbController) createUserDB(username, email string) int {
	var id int
	user, err := dbC.getUserByEmailDB(email)
	if err == nil {
		return user.user_id
	}
	query := "INSERT into users (email, username) VALUES ($1, $2) RETURNING user_id"
	err = dbC.db.QueryRow(query, email, username).Scan(&id)
	log.Print("CreateUser", err, "\n")
	if err != nil {
		return 0
	}
	return id
}

func (dbC *dbController) deleteUserDB(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := dbC.db.Exec(query, userID)
	log.Print("DeleteUser:", err, "\n")
	return err
}

func (dbC *dbController) createFileDB(filename, objectName string, ownerID int) (int, error) {
	var id int
	query := "INSERT into files (file_name, object_name, owner_id) VALUES ($1, $2, $3) RETURNING file_id"
	err := dbC.db.QueryRow(query, filename, objectName, ownerID).Scan(&id)
	log.Print("createFile:", err, "\n")
	return id, err
}

func (dbC *dbController) deleteFileDB(fileID int) error {
	query := "DELETE FROM files WHERE id = $1"
	_, err := dbC.db.Exec(query, fileID)
	log.Print(err, "\n")
	// Delete shares
	query = "DELETE FROM file_shares WHERE file_id = $1"
	_, err2 := dbC.db.Exec(query, fileID)
	log.Print("DeleteFile:", err2, "\n")
	return err
}

func (dbC *dbController) shareFileDB(fileID, userID int) error {
	query := "INSERT into file_shares (file_id, shared_with) VALUES ($1, $2)"
	_, err := dbC.db.Exec(query, fileID, userID)
	log.Print("shareFile:", err, "\n")
	return err
}

func (dbC *dbController) unshareFileDB(fileID, userID int) error {
	query := "DELETE FROM file_shares WHERE file_id = $1 AND shared_with = $2"
	_, err := dbC.db.Exec(query, fileID, userID)
	log.Print("unshareFile: ", err, "\n")
	return err
}

func (dbC *dbController) getUserFilesDB(userID int) ([]fileModel, error) {
	var files []fileModel
	query := "SELECT file_id, object_name, file_name, owner_id, created_at FROM files WHERE owner_id = $1"
	rows, err := dbC.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f fileModel
		err := rows.Scan(&f.file_id, &f.object_name, &f.file_name, &f.owner_id, &f.created_at)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (dbC *dbController) getSharedFilesDB(userID int) ([]fileModel, error) {
	var files []fileModel
	query := "SELECT f.file_id, f.object_name, f.file_name, f.owner_id, f.created_at FROM files f JOIN file_shares s ON f.file_id = s.file_id WHERE s.shared_with = $1"
	rows, err := dbC.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f fileModel
		err := rows.Scan(&f.file_id, &f.object_name, &f.file_name, &f.owner_id, &f.created_at)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (dbC *dbController) getFileSharedFromUserDB(receiverID, sharerID int) ([]fileModel, error) {
	var files []fileModel
	query := "SELECT f.file_id, f.object_name, f.file_name, f.owner_id, f.created_at FROM files f JOIN file_shares s ON f.file_id = s.file_id WHERE s.shared_with = $1 AND f.owner_id = $2"
	rows, err := dbC.db.Query(query, receiverID, sharerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f fileModel
		err := rows.Scan(&f.file_id, &f.object_name, &f.file_name, &f.owner_id, &f.created_at)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func getDb() *dbController {
	var dbC dbController
	dbC.init_db()
	return &dbC
}
