package controllers

import(
	"database/sql"
	"go-file-sharing-platform/utils"
	"go-file-sharing-platform/models"
	"encoding/json"
	"net/http"
)

func UploadFileHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    userID, err := utils.GetUserIDFromToken(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // Assume file upload logic is handled here and you have file details
    file := models.File{
        UserID:     userID,
        FileName:   "example.txt",
        Size:       12345,
        S3URL:      "https://example-bucket.s3.amazonaws.com/example.txt",
        FileType:   "text/plain",
    }

    fileID, err := models.InsertFileMetadata(db, file)
    if err != nil {
        http.Error(w, "Error saving file metadata", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{"file_id": fileID})
}
