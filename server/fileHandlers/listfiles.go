package fileHandlers

import (
	"encoding/json"
	"fileserver/utils"
	"log"
	"net/http"
	"os"
)

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusBadRequest)
		return
	}
	log.Println("RemoteAddress: " + r.RemoteAddr)
	requestedPath, _ := utils.GetFileLocation(r)
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
	json.NewEncoder(w).Encode(fileNames)
}
