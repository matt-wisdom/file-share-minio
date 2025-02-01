package client

type Share struct {
	FileID     int    `json:"file_id"`
	FileSize   int64  `json:"file_size"`
	FileName   string `json:"file_name"`
	SharedAt   string `json:"created_at"`
	SharedWith int    `json:"shared_with"`
}

type SharesResp struct {
	Shares []Share `json:"shares"`
}
