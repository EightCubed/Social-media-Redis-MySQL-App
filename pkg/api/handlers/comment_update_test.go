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

var _ = Describe("CommentPatch", func() {
	var (
		testBody  map[string]interface{}
		router    *mux.Router
		w         *httptest.ResponseRecorder
		r         *http.Request
		jsonBytes []byte
		err       error
	)

	BeforeEach(func() {
		_, router = createFakeSocialMediaHandlerAndRouter()

		// Create a user
		testBody = map[string]interface{}{
			"username": "testuser",
			"email":    "test@example.com",
			"password": "testpassword",
		}

		jsonBytes, err = json.Marshal(testBody)
		Expect(err).ToNot(HaveOccurred())

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusCreated))

		//Create a post
		testBody = map[string]interface{}{
			"user_id": 1,
			"title":   "Post title #1",
			"content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat",
		}

		jsonBytes, err = json.Marshal(testBody)
		Expect(err).ToNot(HaveOccurred())

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/apis/v1/post", bytes.NewBuffer(jsonBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusCreated))

		//Create a comment
		testBody = map[string]interface{}{
			"user_id": 1,
			"content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat",
		}
		jsonBytes, err = json.Marshal(testBody)
		Expect(err).ToNot(HaveOccurred())

		w = httptest.NewRecorder()
		r = httptest.NewRequest("POST", "/apis/v1/post/1/comments", bytes.NewBuffer(jsonBytes))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(http.StatusCreated))
	})

	Describe("POST /apis/v1/post/1/comments", func() {
		Context("when the post exists", func() {
			It("returns a successful response with comment data", func() {
				testBody = map[string]interface{}{
					"content": "Updated Comment",
				}
				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PATCH", "/apis/v1/comments/1", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusOK))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["message"]).To(ContainSubstring("Comment updated successfully"))
			})
		})
		Context("when the comment does not exist", func() {
			It("returns a not found response", func() {
				testBody = map[string]interface{}{
					"content": "Comment 1",
				}
				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PATCH", "/apis/v1/comments/99", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				Expect(w.Body.String()).To(ContainSubstring("Post not found"))
			})
		})
	})
})
