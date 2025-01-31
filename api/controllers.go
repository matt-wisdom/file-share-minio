package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func getFileShare(c *gin.Context) {
	var toUserModel, fromUserModel UserModel
	var err error
	db := getDb()

	fromUser := c.Query("from_user")
	toUser := c.Query("to_user")

	if isEmailAddressRegex((fromUser)) {
		fromUserModel, err = db.getUserByEmailDB(fromUser)
	} else {
		fromUserModel, err = db.getUserByUsernameDB(fromUser)
	}
	if err != nil {
		c.JSON(404, gin.H{"message": fmt.Sprintf("User %s not found", fromUser)})
		return
	}

	if isEmailAddressRegex((toUser)) {
		toUserModel, err = db.getUserByEmailDB(toUser)
	} else {
		toUserModel, err = db.getUserByUsernameDB(toUser)
	}
	if err != nil {
		c.JSON(404, gin.H{"message": fmt.Sprintf("User %s not found", toUser)})
		return
	}

	shares, err := db.getFileSharedFromUserDB(toUserModel.UserID, fromUserModel.UserID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Could not get shares"})
	}
	c.JSON(200, gin.H{"shares": shares})

}

func shareFileUpload(c *gin.Context) {
	db := getDb()
	output := gin.H{"shares": []interface{}{}}
	userEmail := c.PostForm("user_email")
	userName := c.PostForm("user_name")
	toUsers := c.PostFormArray("to_user")
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	userId := db.createUserDB(userName, userEmail)
	for _, file := range files {
		// generate a unique file name
		filename := fmt.Sprintf("%d_%d_%s", userId, time.Now().UnixNano(), file.Filename)
		c.SaveUploadedFile(file, uploadsFolder+filename)
		if err := uploadFileMinio(filename); err != nil {
			c.JSON(500, gin.H{"error": "Failed to upload file"})
			return
		}

		fileId := 0
		for _, toUser := range toUsers {
			toUserModel, err := db.getUserByEmailDB(toUser)
			if err != nil {
				toUserModel, err = db.getUserByUsernameDB(toUser)
				if err != nil {
					if !isEmailAddressRegex(toUser) {
						c.JSON(500, gin.H{"error": "Failed to get user"})
						return
					} else {
						username := strings.Split(toUser, "@")[0]
						userId = db.createUserDB(username, toUser)
						toUserModel, err = db.getUserDB(userId)
						if err != nil {
							c.JSON(500, gin.H{"error": "Failed to create user"})
							return
						}
					}

				}
			}
			if fileId == 0 {
				fileId, err = db.createFileDB(filename, filename, userId)
				if err != nil {
					c.JSON(500, gin.H{"error": "Failed to save file to database"})
					return
				}
			}
			err = db.shareFileDB(fileId, toUserModel.UserID)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to share file"})
				return
			}
			output["shares"] = append(output["shares"].([]interface{}), gin.H{"file_id": fileId, "shared_with": toUserModel.UserID, "shared_at": time.Now().Format("2006-01-02 15:04:05")})
		}

	}
	c.JSON(200, output)
}

func downloadFileResumable(c *gin.Context) {
	db := getDb()
	fileId, err := strconv.Atoi(c.Query("file_id"))
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid file id"})
		return
	}
	var toUser UserModel
	toUserID := c.Query("to_user")
	if isEmailAddressRegex(toUserID) {
		toUser, err = db.getUserByEmailDB(toUserID)
	} else {
		toUser, err = db.getUserByUsernameDB(toUserID)
	}
	receiverID := toUser.UserID
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid receiver id"})
		return
	}

	file, err := db.getFileDB(fileId)
	if err != nil {
		c.JSON(404, gin.H{"message": "File not found"})
		return
	}
	filePath := "downloads/" + file.ObjectName

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		err = downloadFileMinio(file.ObjectName)
	}
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not get object"})
		return
	}
	fileObj, err := os.Open(filePath)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not open file"})
		return
	}
	defer fileObj.Close()

	fileInfo, err := fileObj.Stat()
	if err != nil {
		c.JSON(500, gin.H{"error": "Cannot get file info"})
		return
	}

	fileSize := fileInfo.Size()
	c.Writer.Header().Set("Accept-Ranges", "bytes")
	defer db.setSharedFile(fileId, receiverID)

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		// c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		c.File(filePath)

		return
	}

	byteRange := strings.Split(rangeHeader, "=")
	if len(byteRange) < 2 {
		c.JSON(400, gin.H{"message": "Invalid range header"})
		return
	}

	rangeParts := strings.Split(byteRange[1], "-")
	start, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil || start >= fileSize {
		c.JSON(416, gin.H{"message": "Invalid range start"})
		return
	}

	var end int64 = fileSize - 1
	if len(rangeParts) > 1 && rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil || end >= fileSize {
			end = fileSize - 1
		}
	}

	if start > end {
		c.JSON(416, gin.H{"error": "Invalid range"})
		return
	}

	contentLength := end - start + 1
	if end >= fileSize {
		end = fileSize - 1 // Prevent overflow
		contentLength = end - start + 1
	}

	c.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Writer.Header().Set("Accept-Ranges", "bytes")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Writer.Header().Set("Content-Type", "application/octet-stream") // Adjust for specific file types
	c.Writer.WriteHeader(http.StatusPartialContent)                   // 206 Partial Content

	// Seek to start position
	_, err = fileObj.Seek(start, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to seek file"})
		return
	}

	buffer := make([]byte, 64*1024) // 64 KB buffer
	bytesRemaining := contentLength

	for bytesRemaining > 0 {
		readBytes, err := fileObj.Read(buffer)
		if readBytes > int(bytesRemaining) {
			readBytes = int(bytesRemaining) // Prevent reading beyond range
		}

		if readBytes > 0 {
			c.Writer.Write(buffer[:readBytes])
			c.Writer.Flush()
			bytesRemaining -= int64(readBytes)
		}

		if err == io.EOF {
			break // End of file reached
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "File read error"})
			return
		}
	}

}
