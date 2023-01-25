package multipart_requests

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const DEFAULT_METHOD = "POST"

type MultipartRequest struct {
	// The path where temporary files will be stored on the server
	TempPath string
	// Whether to persist files to the server. By default, files will be removed
	Persist bool
}

// GetFile Gets the file from an incoming form request
func (m *MultipartRequest) GetFile(r *http.Request, name string) (*string, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &handler.Filename, nil
}

// Upload Uploads a local file via an http request as multipart formdata.
func (m *MultipartRequest) Upload(file multipart.File, url, filename string) (*http.Response, error) {
	// create dir if not exist
	err := os.Mkdir("temp", 0755)
	if err != nil {
		log.Printf("Error %s", err.Error())
		return nil, err
	}

	// Create file
	dst, err := os.Create("temp/" + filename)
	defer dst.Close()
	if err != nil {
		return nil, err
	}
	// Copy file over
	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}
	// upload file
	filePath := "temp/" + filename
	tempFile, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("logo", filename)
	if _, err = io.Copy(part, tempFile); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(
		DEFAULT_METHOD,
		url,
		body,
	)

	req.Header.Add("Content-Type", "multipart/form-data")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return res, nil
	}
	// Clean up
	if !m.Persist {
		if err = os.RemoveAll(m.TempPath); err != nil {
			return res, err
		}
	}
	return res, nil
}
