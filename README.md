# Multipart Requests
Http multipart form data requests

### Usage
Create a MultipartRequest type
```go
mr := MultipartRequest{
    TempPath: "temp",
	Persist: false,
}
```

```go
file, err := m.GetFile(r, "image.png")
```

```go
res, err := m.Upload(file, "http://google.com", "image.png", "image")
```



