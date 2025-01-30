package main

type userModel struct {
	user_id  int
	username string
	email    string
}

type fileModel struct {
	file_id     int
	object_name string
	file_name   string
	owner_id    int
	created_at  string
}

type shareModel struct {
	file_id     int
	shared_with int
	shared_at   string
}
