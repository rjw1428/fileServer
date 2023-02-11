package fileHandlers

import (
	"fileserver/utils"
	"log"
	"net/http"
	"os"
)

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusBadRequest)
		return
	}

	requestedPath, fileName := utils.GetFileLocation(r)
	log.Println("Returning file", requestedPath+"/"+fileName)
	file, err := os.Open(requestedPath + "/" + fileName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment: filename="+fileName)
	http.ServeFile(w, r, requestedPath+"/"+fileName)
}
