package handlers

import "net/http"

func (h *SocialMediaHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get user endpoint"))
}
