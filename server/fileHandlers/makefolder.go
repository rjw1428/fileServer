package fileHandlers

import (
	"encoding/json"
	"fileserver/utils"
	"log"
	"net/http"
	"os"
)

func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusBadRequest)
		return
	}

	requestedPath, folderName := utils.GetFileLocation(r)
	log.Println("Createing new folder at", requestedPath)
	err := os.Mkdir(requestedPath+"/"+folderName, os.ModePerm)
	if err != nil {
		http.Error(w, "Directory already exists"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "applicxation/json")
	w.WriteHeader(http.StatusOK)

	resp := &utils.Response{Success: true}
	json.NewEncoder(w).Encode(resp)
}
