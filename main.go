package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const path = "/media/pi/Crucial X6/"
const apiVersion = "v1"
const port = 1428
const MAX_UPLOAD_SIZE = 1024 * 1024 * 1024 * 10

func main() {
	log.Printf("Starting...")
	http.HandleFunc(fmt.Sprintf("/api/%s/files", apiVersion), listFilesHandler)
	http.HandleFunc(fmt.Sprintf("/api/%s/download", apiVersion), downloadFileHandler)
	http.HandleFunc(fmt.Sprintf("/api/%s/upload", apiVersion), uploadFileHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Println("Error starting server:", err)
	}
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	requestedPath, _, err := getFileLocation(r)
	if err != nil {
		http.Error(w, "Unable to read subpath, make sure if there is no subpath being sent, the body contains an empty array", http.StatusBadRequest)
		return
	}
	log.Println("Requesting files from", requestedPath)
	files, err := os.ReadDir(requestedPath)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	type FileInfo struct {
		Name    string `json:"name"`
		IsDir   bool   `json:"is_dir"`
		Size    int64  `json:"size"`
		ModTime int64  `json:"modified"`
	}

	var fileNames []FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info for %s, skipping", file.Name())
		}
		fileNames = append(fileNames, FileInfo{
			Name:    file.Name(),
			Size:    info.Size(),
			IsDir:   file.IsDir(),
			ModTime: info.ModTime().Unix(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((fileNames))
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	requestedFile, fileName, err := getFileLocation(r)
	if err != nil {
		http.Error(w, "Unable to read subpath, make sure if there is no subpath being sent, the body contains an empty array", http.StatusBadRequest)
		return
	}
	log.Println("Returning file", requestedFile)
	file, err := os.Open(requestedFile)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment: filename="+fileName)
	http.ServeFile(w, r, requestedFile)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]
	for _, fileHeader := range files {
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprint("The uploaded file is too big. Please keep file size to less than 10 GB.", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error reading file "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePath := filepath.Join(path, fileHeader.Filename)
		log.Printf("Uploading file %s to %s", fileHeader.Filename, path)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		pr := &Progress{
			TotalSize: fileHeader.Size,
		}
		_, err = io.Copy(out, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully")
}

// UTILITY FUNCTIONS
func getFileLocation(r *http.Request) (string, string, error) {
	var subFolders []string
	err := json.NewDecoder(r.Body).Decode(&subFolders)
	if err != nil {
		return "", "", err
	}
	var fileName string
	if len(subFolders) > 0 {
		fileName = subFolders[len(subFolders)-1]
	}

	requestedFile := path
	for _, subFolder := range subFolders {
		requestedFile = filepath.Join(requestedFile, subFolder)
	}
	return requestedFile, fileName, nil
}

type Progress struct {
	TotalSize int64
	BytesRead int64
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

	fmt.Printf("File upload in progress: %d\n", pr.BytesRead)
}
