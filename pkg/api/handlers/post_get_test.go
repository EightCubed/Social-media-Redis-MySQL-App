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

	Describe("GET /apis/v1/post/{id}", func() {
		Context("when the post exists", func() {
			It("returns a successful response with post data", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/apis/v1/post/1", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusOK))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())

				postBody := responseBody["Post"].(map[string]interface{})
				Expect(postBody["Title"]).To(Equal("Post title #1"))
				Expect(postBody["Content"]).To(Equal("Lorem ipsum dolor sit amet..."))
			})
		})

		Context("when the post get is performed repeatedly", func() {
			var ViewerCount int
			BeforeEach(func() {
				ViewerCount = 10
				for i := 0; i <= ViewerCount; i++ {
					w := httptest.NewRecorder()
					r := httptest.NewRequest("GET", "/apis/v1/post/1", nil)
					r.Header.Set("Content-Type", "application/json")
					router.ServeHTTP(w, r)

					Expect(w.Code).To(Equal(http.StatusOK))
					var responseBody map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &responseBody)
					Expect(err).ToNot(HaveOccurred())

					postBody := responseBody["Post"].(map[string]interface{})
					Expect(postBody["Title"]).To(Equal("Post title #1"))
					Expect(postBody["Content"]).To(Equal("Lorem ipsum dolor sit amet..."))
				}
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

				postBody := responseBody["Post"].(map[string]interface{})
				Expect(postBody["Title"]).To(Equal("Post title #1"))
				Expect(postBody["Content"]).To(Equal("Lorem ipsum dolor sit amet..."))
				Expect(postBody["Views"]).To(Equal(float64(1 + ViewerCount + 1)))
			})
		})

		Context("when the post does not exist", func() {
			It("returns 404", func() {
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/apis/v1/post/99", nil)
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
