package main

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Initialize minio client object.

func initMinio() *minio.Client {
	ctx := context.Background()

	var minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Region: bucketRegion,
	})
	if err != nil {
		panic(err)
	}
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	log.Printf("Make Bucket %v", err)
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s exists\n", bucketName)
		} else {
			panic(err)
		}
	}

	return minioClient
}

var minioClient = initMinio()

func deleteLocalFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete %s: %v\n", filePath, err)
		return err
	}
	log.Printf("Successfully deleted %s\n", filePath)
	return nil

}

func deleteFileMinio(fileName string) error {
	ctx := context.Background()
	err := minioClient.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Failed to delete %s: %v\n", fileName, err)
		return err
	}
	log.Printf("Successfully deleted %s\n", fileName)
	return nil
}

func uploadFileMinio(fileName string) error {
	ctx := context.Background()
	contentType := "application/octet-stream"
	defer deleteLocalFile(uploadsFolder + fileName)
	info, err := minioClient.FPutObject(ctx, bucketName, fileName, uploadsFolder+fileName, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Failed to upload %s: %v\n", fileName, err)
		return err
	}
	log.Printf("Successfully uploaded %s of size %d\n", fileName, info.Size)
	return nil
}

func downloadFileMinio(fileName string) error {
	ctx := context.Background()
	err := minioClient.FGetObject(ctx, bucketName, fileName, "downloads/"+fileName, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Failed to download %s: %v\n", fileName, err)
		return err
	}
	log.Printf("Successfully downloaded %s\n", fileName)
	return nil
}
