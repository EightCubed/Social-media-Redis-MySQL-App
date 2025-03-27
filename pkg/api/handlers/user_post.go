package handlers

import (
	"encoding/json"
	"fmt"
	"go-social-media/pkg/models"
	"go-social-media/pkg/utils"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserPostBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *SocialMediaHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] PostUser handler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	var body UserPostBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user := &models.User{
		Username: body.Username,
		Email:    body.Email,
	}
	log.Printf("[INFO] Creating new user - Username: %s, Email: %s", user.Username, user.Email)
	if err := h.DBWriter.Create(&user).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("[ERROR] Access token generation failed: %v", err)
	}
	refreshToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("[ERROR] Refresh token generation failed: %v", err)
	}

	login := &models.Login{
		PasswordHash: string(hashedPassword),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	log.Printf("[INFO] Creating new login")
	if err := h.DBWriter.Create(&login).Error; err != nil {
		log.Printf("[ERROR] Database insertion failed: %v", err)
		http.Error(w, "Failed to create login", http.StatusInternalServerError)
		return
	}

	cacheKey := fmt.Sprintf("user:%d", user.ID)
	marshalledUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("[ERROR] Marshal error: %v", err)
	}
	err = h.RedisReader.Set(cacheKey, marshalledUser, CACHE_DURATION_LONG).Err()
	if err != nil {
		log.Printf("[ERROR] Cache set error: %v", err)
	}

	log.Printf("[INFO] User created successfully - ID: %d, Username: %s", user.ID, user.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
