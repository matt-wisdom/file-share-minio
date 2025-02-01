package main

type UserModel struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FileModel struct {
	FileID     int    `json:"file_id"`
	ObjectName string `json:"object_name"`
	FileName   string `json:"file_name"`
	OwnerID    int    `json:"owner_id"`
	CreatedAt  string `json:"created_at"`
	FileSize   int64  `json:"file_size"`
}

type ShareModel struct {
	FileID     int    `json:"file_id"`
	SharedWith int    `json:"shared_with"`
	SharedAt   string `json:"shared_at"`
}
