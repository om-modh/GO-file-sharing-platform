package models

import(
	"database/sql"
	"time"
)

type File struct{
	ID         int
    UserID     int
    FileName   string
    UploadDate string
    Size       int64
    S3URL      string
    FileType   string
}

func GetUserFiles(db *sql.DB, userID int) ([]File, error) {
    rows, err := db.Query(`
        SELECT id, user_id, file_name, upload_date, size, s3_url, file_type 
        FROM files 
        WHERE user_id = $1
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var files []File
    for rows.Next() {
        var file File
        if err := rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.UploadDate, &file.Size, &file.S3URL, &file.FileType); err != nil {
            return nil, err
        }
        files = append(files, file)
    }
    return files, nil
}

func InsertFileMetadata(db *sql.DB, file File) (int, error) {
    var fileID int
    err := db.QueryRow(`
        INSERT INTO files (user_id, file_name, size, s3_url, file_type) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id
    `, file.UserID, file.FileName, file.Size, file.S3URL, file.FileType).Scan(&fileID)
    if err != nil {
        return 0, err
    }
    return fileID, nil
}

func GetExpiredFiles(db *sql.DB, expirationTime time.Time) ([]File, error) {
    query := `
        SELECT id, user_id, file_name, size, s3_url, file_type, created_at
        FROM files
        WHERE created_at < $1
    `
    
    rows, err := db.Query(query, expirationTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var files []File
    for rows.Next() {
        var file File
        if err := rows.Scan(&file.ID, &file.UserID, &file.FileName, &file.Size, &file.S3URL, &file.FileType); err != nil {
            return nil, err
        }
        files = append(files, file)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return files, nil
}

// DeleteFile removes a file record from the database
func DeleteFile(db *sql.DB, fileID int) error {
    query := `DELETE FROM files WHERE id = $1`
    _, err := db.Exec(query, fileID)
    return err
}