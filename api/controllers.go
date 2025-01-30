package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func shareFileUpload(c *gin.Context) {
	// Get file and file content
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
			err = db.shareFileDB(fileId, toUserModel.user_id)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to share file"})
				return
			}
			output["shares"] = append(output["shares"].([]interface{}), gin.H{"file_id": fileId, "shared_with": toUserModel.user_id, "shared_at": time.Now().Format("2006-01-02 15:04:05")})
		}

	}
	c.JSON(200, output)
	return
}
