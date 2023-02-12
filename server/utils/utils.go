package utils

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type Progress struct {
	TotalSize int64
	BytesRead int64
}

type Response struct {
	Success bool `json:"success"`
}

func GetFileLocation(r *http.Request) (string, string) {
	subPath := r.URL.Query().Get("path")
	fileName := r.URL.Query().Get("file")
	requestedFile := filepath.Join(ROOT_DIR, subPath)
	return requestedFile, fileName
}

func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	pr.Print()
	return
}

// Print displays the current progress of the file upload
// each time Write is called
func (pr *Progress) Print() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("DONE!")
		return
	}
	var percent = float32(pr.BytesRead) / float32(pr.TotalSize) * 100
	fmt.Printf("File upload in progress: %.2f %%\n", percent)
}
