# 🚀 Go File Sharing API

## 📌 Overview
This project is an extremely simple **file-sharing API** built using **Gin Gonic**, **PostgreSQL**, and **MinIO** for storage. It supports **file uploads, sharing, and resumable downloads**.

## 🏗️ Tech Stack
- **Go** (Gin Gonic framework)
- **PostgreSQL** (Database for storing file metadata)
- **MinIO** (Object storage for files)
- **Docker & Docker Compose** (Containerized environment)

## 📂 Features
- 📤 **Upload files**
- 🔗 **Share files** with other users
- 📥 **Resumable downloads** with `Range` headers
- 📦 **Dockerized setup** for easy deployment

## 🏃‍♂️ Quick Start

### 🔹 1. Clone the Repository
```sh
git clone https://github.com/your-username/go-file-sharing-api.git
cd go-file-sharing-api
```

### 🔹 2. Set Up Environment Variables
Create a `.env` file inside the `api/` directory and configure:
```env
POSTGRES_USER=fileshare
POSTGRES_PASSWORD="123456"
POSTGRES_DB=fileshare
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

MINIO_ENDPOINT=minio:9000
MINIO_ROOT_USER=rootuser
MINIO_ROOT_PASSWORD=rootpass
MINIO_ID=rootuser
MINIO_SK=rootpass

# Optional values (remove if not needed)
MINIO_DISABLE_SSL=true
MINIO_BUCKET=fileshare
MINIO_REGION=us-east-1
```

### 🔹 3. Start Services with Docker Compose
```sh
docker-compose up -d --build
```

### 🔹 4. API Endpoints
| Method | Endpoint | Description |
|--------|---------|-------------|
| `POST` | `/share-file` | Uploads a file |
| `GET` | `/shares` | Get files shared from a user |
| `GET`  | `/download` | Downloads a file (supports range requests) |



## 🔄 Resumable Download Example
```sh
curl -H "Range: bytes=0-1023" -o part1.bin "http://localhost:8080/download?file_id=file123&receiver_id=123"
```

## 🛠️ Development
To run the API locally without Docker:
```sh
cd api
go mod tidy
go run main.go
```

## 🔥 Contributing
Pull requests are welcome! Please open an issue first to discuss changes.

## 📜 License
MIT License. See `LICENSE` for details.

