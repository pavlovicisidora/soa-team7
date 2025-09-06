package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024) // 10 MB

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "The uploaded file is too big.", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		http.Error(w, "Image key not found in form.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Error creating upload directory: %v", err)
		http.Error(w, "Could not create upload directory.", http.StatusInternalServerError)
		return
	}

	uniqueFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))
	filePath := filepath.Join(uploadDir, uniqueFileName)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating the file on server.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving the file.", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully uploaded file: %s", uniqueFileName)

	publicPath := "/uploads/" + uniqueFileName

	response := map[string]string{"filePath": publicPath}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
