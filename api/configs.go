package main

import "os"

var endpoint = os.Getenv("MINIO_ENDPOINT")

var accessKeyID = os.Getenv("MINIO_ID")
var secretAccessKey = os.Getenv("MINIO_SK")
var useSSL = os.Getenv("MINIO_DISABLE_SSL") != "true"
var bucketName = os.Getenv("MINIO_BUCKET")
var bucketRegion = os.Getenv("MINIO_REGION")
var uploadsFolder = "uploads/"

var postgresUser = os.Getenv("POSTGRES_USER")
var postgresPassword = os.Getenv("POSTGRES_PASSWORD")
var postgresHost = os.Getenv("POSTGRES_HOST")
var postgresPort = os.Getenv("POSTGRES_PORT")
var postgresDB = os.Getenv("POSTGRES_DB")
