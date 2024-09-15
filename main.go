package main

import(
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-file-sharing-platform/routes"
	"github.com/joho/godotenv"
	_"github.com/lib/pq"
)

func main(){
	err:=godotenv.Load()
	if err!=nil {
		log.Fatal("Error Loading .env file")
	}

	dbHost:= os.Getenv("DB_HOST")
	dbPort:= os.Getenv("DB_PORT")
	dbUser:= os.Getenv("DB_USER")
	dbPassword:= os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	pqsqlInfo:= fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db,err:=sql.Open("postgres", pqsqlInfo)
	if err!=nil {
		log.Fatal(err)
	}
	defer db.Close()

	routes.AuthRoutes(db)
	http.HandleFunc("/upload", routes.UploadFileHandler)

	http.HandleFunc("/files", routes.GetFileHandler)
	http.HandleFunc("/share/:file_id", routes.ShareFilehandler)
	http.HandleFunc("/files/search", routes.SearchFilesHandler)  // Search files endpoint


	log.Println("Srever running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}