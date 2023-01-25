package multipart_requests

import (
	"bytes"
	"fmt"
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
	// Set to true to get error logs. Default is false
	Debug bool
	// The request url
	Url string
}

// GetFile Gets the file from an incoming form request
func (m *MultipartRequest) GetFile(r *http.Request, name string) (*string, multipart.File, error) {
	file, handler, err := r.FormFile(name)
	if err != nil {
		return nil, nil, err
	}
	//defer file.Close()
	return &handler.Filename, file, nil
}

// Upload a local file via an http request as multipart formdata.
func (m *MultipartRequest) Upload(file multipart.File, filename, field string) (*http.Response, error) {
	var err error
	if m.Url == "" {
		panic("You must provide a Url value")
	}
	if m.TempPath == "" {
		m.TempPath = "temp"
	}
	// create dir if not exist
	err = os.Mkdir("temp", 0755)
	if err != nil {
		if m.Debug {
			log.Printf("Error creating temprary file: %s", err.Error())
		}
	}

	// Create file
	filePath := fmt.Sprintf("%s/%s", m.TempPath, filename)
	dst, err := os.Create(filePath)
	defer dst.Close()
	if err != nil {
		return nil, err
	}
	// Copy file over
	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}
	// upload file
	tempFile, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(field, filename)
	if _, err = io.Copy(part, tempFile); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(
		DEFAULT_METHOD,
		m.Url,
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
