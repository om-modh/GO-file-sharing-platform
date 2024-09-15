package routes

import (
	"database/sql"
	"encoding/json"
	"go-file-sharing-platform/utils"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type FileMetadata struct {
	ID         int    `json:"id"`
    FileName   string `json:"file_name"`
    FileSize   int64  `json:"file_size"`
    UploadDate string `json:"upload_date"`
    FilePath   string `json:"file_path"`
}

func GetFileHandler(w http.ResponseWriter, r *http.Request){
	userID, err := utils.GetUserIDFromToken(r)
    if err != nil {
        http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
        return
    }

	db,_ := sql.Open("postres", os.Getenv("DB_URL"))
	defer db.Close()

	var files []FileMetadata

	rows, err:= db.Query(`SELECT id, file_name, file_size, upload_date, file_path FROM files WHERE user_id=$1`, userID)

	if err!=nil {
		http.Error(w, "Error while file Retrieving", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var file FileMetadata
		err:= rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.UploadDate, &file.FilePath)
		if err!=nil {
			http.Error(w, "Error reading file data", http.StatusInternalServerError)
			return 
		}
		files = append(files, file)
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(files)
}


func ShareFilehandler(w http.ResponseWriter, r *http.Request){
	vars:= mux.Vars(r)
	fileID:=vars["file_id"]

	db,_ := sql.Open("postgres", os.Getenv("DB_URL"))
	defer db.Close()

	var file FileMetadata
	err:= db.QueryRow(`SELECT file_path FROM files WHERE id=$1`, fileID).Scan(&file.FilePath)
	if err!=nil {
		http.Error(w, "File Not Found", http.StatusNotFound)
		return
	}

	publicURL := "http://localhost:8080/uploads/" + file.FilePath
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url":publicURL})
}
