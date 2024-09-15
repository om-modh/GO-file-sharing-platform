package routes

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadFileHandler: This function handles file uploads
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the incoming file from the user's request
    r.ParseMultipartForm(10 << 20) // Limit file size to 10MB
    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error receiving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Save the file locally in a folder called 'uploads'
    localPath := saveFileLocally(file, handler.Filename)

    // Save file metadata in the database
    storeFileMetadata(handler.Filename, localPath, handler.Size)

    fmt.Fprintf(w, "File uploaded successfully: %s\n", localPath)
}


// saveFileLocally: This saves the file to a local folder called 'uploads'
func saveFileLocally(file multipart.File, fileName string) string {
    // Create the 'uploads' folder if it doesn't exist
    os.MkdirAll("uploads", os.ModePerm)

    // Define the path where the file will be saved
    fullPath := filepath.Join("uploads", fmt.Sprintf("%d_%s", time.Now().Unix(), fileName))

    // Create the file on the local system
    outFile, _ := os.Create(fullPath)
    defer outFile.Close()

    // Copy the uploaded file's content to the local file
    io.Copy(outFile, file)

    return fullPath
}

// storeFileMetadata: This stores file information (metadata) in the database
func storeFileMetadata(fileName string, filePath string, size int64) {
    db, _ := sql.Open("postgres", os.Getenv("DB_URL"))
    defer db.Close()

    // Save the metadata (file info) into the 'files' table
    query := `INSERT INTO files (user_id, file_name, file_path, file_size, upload_date) 
              VALUES (1, $1, $2, $3, $4)` // Example: user_id is hardcoded for simplicity
    db.Exec(query, fileName, filePath, size, time.Now())
}
