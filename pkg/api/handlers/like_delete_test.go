package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LikeDelete", func() {
	var router *mux.Router

	BeforeEach(func() {
		_, router = createFakeSocialMediaHandlerAndRouter()

		// Create users
		for i := range 2 {
			body := map[string]interface{}{
				"username": "testuser" + string(rune(i)),
				"email":    "test" + string(rune(i)) + "@example.com",
				"password": "testpassword",
			}
			jsonBytes, _ := json.Marshal(body)
			req := httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			Expect(res.Code).To(Equal(http.StatusCreated))
		}

		// Create a post by user 1
		postBody := map[string]interface{}{
			"user_id": 1,
			"title":   "Sample Title",
			"content": "Sample Content",
		}
		jsonBytes, _ := json.Marshal(postBody)
		req := httptest.NewRequest("POST", "/apis/v1/post", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		Expect(res.Code).To(Equal(http.StatusCreated))

		// Like post 1 by user 1 and 2
		for i := 1; i <= 2; i++ {
			url := "/apis/v1/post/1/likes?user_id=" + string(rune('0'+i))
			req := httptest.NewRequest("POST", url, nil)
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			Expect(res.Code).To(Equal(http.StatusCreated))
		}
	})

	Describe("DELETE /apis/v1/post/{id}/likes", func() {
		Context("When the like exists", func() {
			It("successfully deletes the like and returns updated like count", func() {
				req := httptest.NewRequest("DELETE", "/apis/v1/post/1/likes?user_id=1", nil)
				req.Header.Set("Content-Type", "application/json")
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				Expect(res.Code).To(Equal(http.StatusOK))
				var resp map[string]interface{}
				err := json.Unmarshal(res.Body.Bytes(), &resp)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp["message"]).To(ContainSubstring("Like deleted successfully"))

				// Confirm like count decreased
				getReq := httptest.NewRequest("GET", "/apis/v1/post/1", nil)
				getRes := httptest.NewRecorder()
				router.ServeHTTP(getRes, getReq)

				Expect(getRes.Code).To(Equal(http.StatusOK))
				var post map[string]interface{}
				err = json.Unmarshal(getRes.Body.Bytes(), &post)
				Expect(err).ToNot(HaveOccurred())
				Expect(post["NumberOfLikes"]).To(Equal(float64(1)))
			})
		})

		Context("When the user does not exist", func() {
			It("returns a Not Found error", func() {
				req := httptest.NewRequest("DELETE", "/apis/v1/post/1/likes?user_id=99", nil)
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)

				Expect(res.Code).To(Equal(http.StatusNotFound))
				Expect(res.Body.String()).To(ContainSubstring("User not found"))
			})
		})

		Context("When the like does not exist", func() {
			It("returns a Not Found error", func() {
				// First delete the like
				req := httptest.NewRequest("DELETE", "/apis/v1/post/1/likes?user_id=2", nil)
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)
				Expect(res.Code).To(Equal(http.StatusOK))

				// Delete again - should fail
				req = httptest.NewRequest("DELETE", "/apis/v1/post/1/likes?user_id=2", nil)
				res = httptest.NewRecorder()
				router.ServeHTTP(res, req)
				Expect(res.Code).To(Equal(http.StatusNotFound))
				Expect(res.Body.String()).To(ContainSubstring("Like not found"))
			})
		})
	})
})
