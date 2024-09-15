// controllers/authController.go
package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "go-file-sharing-platform/models"
    "go-file-sharing-platform/utils"
    
)

func Register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    passwordHash, err := utils.HashPassword(input.Password)
    if err != nil {
        http.Error(w, "Error creating password", http.StatusInternalServerError)
        return
    }

    err = models.CreateUser(db, input.Email, passwordHash)
    if err != nil {
        http.Error(w, "Error saving user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// Login handles user login, generates a JWT token, and responds with the token
func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Decode the request body
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Find the user by email
    user, err := models.FindUserByEmail(db, input.Email)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Check if the password is correct
    if err := utils.CheckPasswordHash(input.Password, user.PasswordHash); err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    // Generate a JWT token for the user
    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    // Respond with the token
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}
