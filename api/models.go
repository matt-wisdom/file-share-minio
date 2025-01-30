package main

type UserModel struct {
	user_id  int
	username string
	email    string
}

type FileModel struct {
	file_id     int
	object_name string
	file_name   string
	owner_id    int
	created_at  string
}

type ShareModel struct {
	file_id     int
	shared_with int
	shared_at   string
}
