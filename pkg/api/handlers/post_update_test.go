package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("PostUpdate", func() {
	var (
		testBody  map[string]interface{}
		router    *mux.Router
		w         *httptest.ResponseRecorder
		r         *http.Request
		jsonBytes []byte
		err       error
	)

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

	// Happy path
	Context("when request body is valid", func() {
		BeforeEach(func() {
			testBody = map[string]interface{}{
				"title":   "Update title #1",
				"content": "Updated content Ut enim ad minim veniam, ullamco laboris nisi ut aliquip ex ea commodo consequat",
			}
			jsonBytes, err = json.Marshal(testBody)
			Expect(err).ToNot(HaveOccurred())

			w = httptest.NewRecorder()
			r = httptest.NewRequest("PATCH", "/apis/v1/post/1", bytes.NewBuffer(jsonBytes))
			r.Header.Set("Content-Type", "application/json")
		})

		It("should handle valid JSON body and return success", func() {
			router.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusOK))
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			Expect(err).ToNot(HaveOccurred())
			fmt.Println("responseBody", responseBody)
			Expect(responseBody["post_id"]).To(Equal(float64(1)))
			Expect(responseBody["message"]).To(Equal("Post updated successfully"))
		})
	})

	Context("when request body is not valid", func() {
		BeforeEach(func() {
			testBody = map[string]interface{}{
				"title":   "Update title #1",
				"content": "Updated content Ut enim ad minim veniam, ullamco laboris nisi ut aliquip ex ea commodo consequat",
				"extra":   "field",
			}
			jsonBytes, err = json.Marshal(testBody)
			Expect(err).ToNot(HaveOccurred())

			w = httptest.NewRecorder()
			r = httptest.NewRequest("PATCH", "/apis/v1/post/1", bytes.NewBuffer(jsonBytes))
			r.Header.Set("Content-Type", "application/json")
		})

		It("should return error", func() {
			router.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when post does not exist", func() {
		BeforeEach(func() {
			testBody = map[string]interface{}{
				"title":   "Update title #1",
				"content": "Updated content Ut enim ad minim veniam, ullamco laboris nisi ut aliquip ex ea commodo consequat",
			}
			jsonBytes, err = json.Marshal(testBody)
			Expect(err).ToNot(HaveOccurred())

			w = httptest.NewRecorder()
			r = httptest.NewRequest("PATCH", "/apis/v1/post/2", bytes.NewBuffer(jsonBytes))
			r.Header.Set("Content-Type", "application/json")
		})

		It("should return error", func() {
			router.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})
})
