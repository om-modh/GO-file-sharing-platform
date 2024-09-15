// routes/authRoutes.go
package routes

import (
    "database/sql"
    "go-file-sharing-platform/controllers"
    "net/http"
)

func AuthRoutes(db *sql.DB) {
    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            controllers.Register(db, w, r)
        }
    })

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            controllers.Login(db, w, r)
        }
    })
}
