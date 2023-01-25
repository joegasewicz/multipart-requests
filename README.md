# Multipart Requests
Http multipart form data requests

### Install
```
go get -u github.com/joegasewicz/multipart-requests
```

### Usage
Upload the file as multipart form data within a http request
```go
m := multipart_requests.MultipartRequest{
    TempPath: "uploads", // Default is temp
    Persist:  true, // Default is false
	Url:      "http://127.0.0.1:8000/example", // Required
    Debug:    true, // Default is false
}
// If a file is being obtained via a form request then use this helper
fileName, file, err := m.GetFile(r, "logo")
// Upload the file as multipart form data within a http request
_, err = m.Upload(file,  *fileName, "logo")
```

