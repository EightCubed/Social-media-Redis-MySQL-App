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

var _ = Describe("PostPost", func() {
	var (
		router *mux.Router
	)

	BeforeEach(func() {
		_, router = createFakeSocialMediaHandlerAndRouter()

		// Create user
		testBody := map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "testpassword",
		}
		jsonBytes, _ := json.Marshal(testBody)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusCreated))

		testBody = map[string]interface{}{
			"username": "testuser_1",
			"email":    "test_1@example.com",
			"password": "testpassword",
		}
		jsonBytes, _ = json.Marshal(testBody)
		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusCreated))

		// Create post
		testBody = map[string]interface{}{
			"user_id": 1,
			"title":   "Post title #1",
			"content": "Lorem ipsum dolor sit amet...",
		}
		jsonBytes, _ = json.Marshal(testBody)
		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/apis/v1/post", bytes.NewBuffer(jsonBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusCreated))
	})

	Describe("POST /apis/v1/post/{id}/likes", func() {
		Context("When the post exists", func() {
			BeforeEach(func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/apis/v1/post/1/likes?user_id=1", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusCreated))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["message"]).To(ContainSubstring("Like added successfully"))
			})
			It("returns a successful response with post data", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/apis/v1/post/1", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusOK))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["NumberOfLikes"]).To(Equal(float64(1)))
			})

			It("returns a successful response with post data", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/apis/v1/post/1/likes?user_id=2", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusCreated))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["message"]).To(ContainSubstring("Like added successfully"))

				w = httptest.NewRecorder()
				r = httptest.NewRequest("GET", "/apis/v1/post/1", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusOK))
				err = json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["NumberOfLikes"]).To(Equal(float64(2)))
			})
		})
		Context("When the post exists", func() {
			It("returns a bad request", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/apis/v1/post/1/likes?user_id=99", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
