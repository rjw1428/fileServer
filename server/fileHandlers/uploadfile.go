package fileHandlers

import (
	"encoding/json"
	"fileserver/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, utils.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(utils.MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["file"]
	path := r.MultipartForm.Value["path"][0]
	for _, fileHeader := range files {
		if fileHeader.Size > utils.MAX_UPLOAD_SIZE {
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

		filePath := filepath.Join(utils.ROOT_DIR, path, fileHeader.Filename)
		log.Printf("Uploading file %s to %s", fileHeader.Filename, utils.ROOT_DIR)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		pr := &utils.Progress{
			TotalSize: fileHeader.Size,
		}
		_, err = io.Copy(out, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := &utils.Response{Success: true}
	json.NewEncoder(w).Encode(resp)
}
