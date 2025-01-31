package main

import (
	"context"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

const (
	testBucket   = "test-bucket"
	testFile     = "testfile.txt"
	testContent  = "Hello, MinIO!"
	testDownload = "downloads/testfile.txt"
)

func setupTestMinio() *minio.Client {
	ctx := context.Background()

	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	client.MakeBucket(ctx, testBucket, minio.MakeBucketOptions{})
	return client
}

func TestUploadFileMinio(t *testing.T) {
	// Setup MinIO
	minioClient = setupTestMinio()
	bucketName = testBucket
	uploadsFolder = "./"

	// Create a test file
	os.WriteFile(testFile, []byte(testContent), 0644)
	defer os.Remove(testFile)

	err := uploadFileMinio(testFile)
	assert.Nil(t, err)

	// Check if file exists in MinIO
	ctx := context.Background()
	_, err = minioClient.StatObject(ctx, testBucket, testFile, minio.StatObjectOptions{})
	assert.Nil(t, err)
}

func TestDownloadFileMinio(t *testing.T) {
	// Setup MinIO
	minioClient = setupTestMinio()
	bucketName = testBucket

	// Upload a test file
	os.WriteFile(testFile, []byte(testContent), 0644)
	defer os.Remove(testFile)
	minioClient.FPutObject(context.Background(), testBucket, testFile, testFile, minio.PutObjectOptions{})

	// Test download
	err := downloadFileMinio(testFile)
	assert.Nil(t, err)

	// Check if the file is downloaded
	data, err := os.ReadFile(testDownload)
	assert.Nil(t, err)
	assert.Equal(t, testContent, string(data))

	defer os.Remove(testDownload)
}

func TestDeleteFileMinio(t *testing.T) {
	// Setup MinIO
	minioClient = setupTestMinio()
	bucketName = testBucket

	// Upload a test file
	os.WriteFile(testFile, []byte(testContent), 0644)
	defer os.Remove(testFile)
	minioClient.FPutObject(context.Background(), testBucket, testFile, testFile, minio.PutObjectOptions{})

	// Delete file from MinIO
	err := deleteFileMinio(testFile)
	assert.Nil(t, err)

	// Check if file is deleted
	ctx := context.Background()
	_, err = minioClient.StatObject(ctx, testBucket, testFile, minio.StatObjectOptions{})
	assert.NotNil(t, err)
}
