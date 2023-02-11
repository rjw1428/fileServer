package main

import (
	"fileserver/fileHandlers"
	"fileserver/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	log.Printf("Starting...")

	mux := http.NewServeMux()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	mux.HandleFunc(fmt.Sprintf("/api/%s/files", utils.API_VERSION), fileHandlers.ListFilesHandler)
	mux.HandleFunc(fmt.Sprintf("/api/%s/download", utils.API_VERSION), fileHandlers.DownloadFileHandler)
	mux.HandleFunc(fmt.Sprintf("/api/%s/upload", utils.API_VERSION), fileHandlers.UploadFileHandler)

	handler := cors.Handler(mux)

	err := http.ListenAndServe(fmt.Sprintf(":%d", utils.PORT), handler)
	if err != nil {
		log.Println("Error starting server:", err)
	}

}
