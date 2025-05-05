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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		log.Printf("[ERROR] Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Password hashing failed: %v", err)
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	login := &models.Login{
		PasswordHash: string(hashedPassword),
	}
	log.Printf("[INFO] Creating new login")
	if err := h.DBWriter.Create(login).Error; err != nil {
		log.Printf("[ERROR] Login creation failed: %v", err)
		http.Error(w, "Failed to create login", http.StatusInternalServerError)
		return
	}

	// Create user with login ID
	user := &models.User{
		Username: body.Username,
		Email:    body.Email,
		LoginID:  login.ID,
	}
	log.Printf("[INFO] Creating new user - Username: %s, Email: %s", user.Username, user.Email)
	if err := h.DBWriter.Create(user).Error; err != nil {
		log.Printf("[ERROR] User creation failed: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("[ERROR] Access token generation failed: %v", err)
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("[ERROR] Refresh token generation failed: %v", err)
	}

	login.AccessToken = accessToken
	login.RefreshToken = refreshToken
	if err := h.DBWriter.Save(login).Error; err != nil {
		log.Printf("[ERROR] Failed to update login with tokens: %v", err)
	}

	cacheKey := fmt.Sprintf("user:%d", user.ID)
	marshalledUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("[WARN] Marshal error: %v", err)
	} else {
		if err := h.RedisReader.Set(cacheKey, marshalledUser, CACHE_DURATION_LONG).Err(); err != nil {
			log.Printf("[WARN] Cache set error: %v", err)
		}
	}

	log.Printf("[INFO] User created successfully - ID: %d, Username: %s", user.ID, user.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
