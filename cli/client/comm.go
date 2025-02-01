package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

// FileShareServer represents a file-sharing server with a download location.
type FileShareServer struct {
	ServerAddress  string // Server's base URL
	DownloadFolder string // Directory where downloaded files will be stored
}

// downloadFile downloads a file in chunks using HTTP Range requests.
func (fs *FileShareServer) downloadFile(fileID int, fileSize int64, toUser, filename string) error {
	// Construct the download URL with query parameters
	parsedUrl, _ := url.Parse(fs.ServerAddress + "/download")
	params := url.Values{}
	fId := strconv.Itoa(fileID)
	params.Add("file_id", fId)
	params.Add("to_user", toUser)
	parsedUrl.RawQuery = params.Encode()

	client := &http.Client{}

	// Open the destination file for writing
	outFile, err := os.OpenFile(fs.DownloadFolder+"/"+filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Define chunk size (20MB) and calculate the number of chunks needed
	var chunkSize int64 = 1024 * 1024 * 20
	numChunks := 1 + ((fileSize + chunkSize - 1) / chunkSize)

	// Use WaitGroup to manage concurrent chunk downloads
	var wg sync.WaitGroup
	errChan := make(chan error, numChunks)

	// Download each chunk concurrently
	for i := 0; int64(i) < numChunks; i++ {
		offset := int64(i) * chunkSize
		end := offset + chunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}

		wg.Add(1)
		go func(offset, end int64) {
			defer wg.Done()

			// Create a GET request with the appropriate byte range
			req, err := http.NewRequestWithContext(context.Background(), "GET", parsedUrl.String(), nil)
			if err != nil {
				errChan <- err
				return
			}

			// Set Range header to request specific byte range
			rangeHeader := fmt.Sprintf("bytes=%d-%d", offset, end)
			req.Header.Set("Range", rangeHeader)

			// Perform the request
			resp, err := client.Do(req)
			if err != nil {
				errChan <- err
				return
			}
			defer resp.Body.Close()

			// Ensure server supports partial content response
			if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)

				errChan <- fmt.Errorf("unexpected status code %d for chunk [%d-%d] %s", resp.StatusCode, offset, end, string(body))
				return
			}

			// Read the response body
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				errChan <- err
				return
			}

			// Write data to file at the correct offset
			_, err = outFile.WriteAt(data, offset)
			if err != nil {
				errChan <- err
				return
			}

			fmt.Printf("Downloaded chunk [%d-%d] (%d bytes)\n", offset, end, len(data))

		}(offset, end)
	}

	// wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for errors in any chunk download
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	fmt.Println("File downloaded successfully")
	return nil
}

// ReceiveFiles retrieves shared files between users from the server.
func (fs *FileShareServer) ReceiveFiles(fromUser, toUser string) error {
	// Construct the API URL with parameters
	parsedUrl, err := url.Parse(fs.ServerAddress + "/shares")
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("to_user", toUser)
	params.Add("from_user", fromUser)
	parsedUrl.RawQuery = params.Encode()

	// Create a GET request to fetch shared files
	req, err := http.NewRequestWithContext(context.Background(), "GET", parsedUrl.String(), nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body
	var shares SharesResp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Parse JSON response
	err = json.Unmarshal(body, &shares)
	if err != nil {
		return err
	}

	// Print the received files
	fmt.Println(shares.Shares)
	fmt.Printf("Got %d files from %s\n", len(shares.Shares), fromUser)
	for _, share := range shares.Shares {
		fmt.Printf("--> %s - %d received: %s\n", share.FileName, share.FileSize, share.SharedAt)
	}
	for _, share := range shares.Shares {
		fmt.Printf("--> Downloading %s", share.FileName)
		err = fs.downloadFile(share.FileID, share.FileSize, toUser, share.FileName)
		if err != nil {
			fmt.Printf("Error downloading %s: %v", share.FileName, err)
		}
	}
	return nil
}

// ShareFiles uploads files from one user to another.
func (fs *FileShareServer) ShareFiles(files []string, fromUsername, fromUserEmail, toUser string) error {
	var requestBody bytes.Buffer

	// Verify all files exist before proceeding
	for _, file := range files {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			return err
		}
	}

	// Create multipart writer for file upload
	writer := multipart.NewWriter(&requestBody)

	// Add metadata fields
	_ = writer.WriteField("user_name", fromUsername)
	_ = writer.WriteField("user_email", fromUserEmail)
	_ = writer.WriteField("to_user", toUser)

	// Attach each file to the request
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a form file field
		part, err := writer.CreateFormFile("files[]", file.Name())
		if err != nil {
			return err
		}

		// Copy file contents into the form
		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}

	// Close the multipart writer
	if err := writer.Close(); err != nil {
		return err
	}

	// Create an HTTP POST request to send the file
	request, err := http.NewRequestWithContext(context.Background(), "POST", fs.ServerAddress+"/share-file", &requestBody)
	if err != nil {
		return err
	}

	// Set the appropriate content type
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and display the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Body: ", string(respBody))
	return nil
}
