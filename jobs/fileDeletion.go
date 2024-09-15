package jobs

import(
	"database/sql"
	"log"
	"os"
	"go-file-sharing-platform/models"
	"time"
)

const (
	fileRetentionPeriod = 24*time.Hour
)

func RunFileDeletionJob(db *sql.DB){
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select{
		case<-ticker.C:
			deleteExpiredFiles(db)
		}
	}
}

func deleteExpiredFiles(db *sql.DB){
	log.Println("Checking for expired files...")

    now := time.Now()
    expirationTime := now.Add(-fileRetentionPeriod)

    files, err := models.GetExpiredFiles(db, expirationTime)
    if err != nil {
        log.Println("Error fetching expired files:", err)
        return
    }

    for _, file := range files {
        err := os.Remove("./uploads/" + file.FileName) // Delete from local storage
        if err != nil {
            log.Println("Error deleting file from local storage:", err)
        }

        err = models.DeleteFile(db, file.ID) // Delete from database
        if err != nil {
            log.Println("Error deleting file from database:", err)
        }
    }
}
