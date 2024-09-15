package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"go-file-sharing-platform/utils"
	_"github.com/lib/pq"
)

type FileSearchParams struct {
    FileName   string `json:"file_name"`
    UploadDate string `json:"upload_date"`
    FileType   string `json:"file_type"`
}

func SearchFilesHandler(w http.ResponseWriter, r *http.Request){
	userID, err:= utils.GetUserIDFromToken(r)
	if err!=nil {
		http.Error(w, "Unauthorized: " + err.Error(), http.StatusUnauthorized)
		return 
	}

	var params FileSearchParams
	if err:= json.NewDecoder(r.Body).Decode(&params); err!=nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

	query := "SELECT id, file_name, file_size, upload_date, file_path FROM files WHERE user_id = $1"
    args := []interface{}{userID}

	if params.FileName != "" {
        query += " AND file_name ILIKE $2"
        args = append(args, "%"+params.FileName+"%")
    }
    if params.UploadDate != "" {
        query += " AND upload_date = $3"
        args = append(args, params.UploadDate)
    }
    if params.FileType != "" {
        query += " AND file_path ILIKE $4"
        args = append(args, "%"+params.FileType+"%")
    }

	rows, err := db.Query(query, args...)
    if err != nil {
        http.Error(w, "Error retrieving files", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

	var files []FileMetadata
    for rows.Next() {
        var file FileMetadata
        err := rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.UploadDate, &file.FilePath)
        if err != nil {
            http.Error(w, "Error processing file data", http.StatusInternalServerError)
            return
        }
        files = append(files, file)
    }

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(files)
}