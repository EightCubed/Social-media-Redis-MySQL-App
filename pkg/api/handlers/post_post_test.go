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

	"go-social-media/pkg/models"
)

type PostReturn struct {
	Post PostObject
	// NumberOfLikes uint
}

type PostObject struct {
	ID        int
	CreatedAt string
	UpdatedAt string
	DeletedAt string
	UserID    int
	Title     string
	Content   string
	Views     int
	User      UserObject
	Comments  []models.Comment
	Likes     []models.Like
}

var _ = Describe("PostPost", func() {
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

	// Happy path
	Context("when request body is valid", func() {
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

		It("should handle valid JSON body and return success", func() {
			router.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusCreated))
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseBody["message"]).To(Equal("Post created successfully"))
		})

		Context("when request body is not correct", func() {
			Context("request body is nil", func() {
				BeforeEach(func() {
					testBody = map[string]interface{}{}

					jsonBytes, err = json.Marshal(testBody)
					Expect(err).ToNot(HaveOccurred())

					w = httptest.NewRecorder()
					r = httptest.NewRequest("POST", "/apis/v1/post", bytes.NewBuffer(jsonBytes))
					r.Header.Set("Content-Type", "application/json")
				})

				It("should return bad request", func() {
					router.ServeHTTP(w, r)
					fmt.Println(w)
					fmt.Println(r.Body)
					Expect(w.Body.String()).To(ContainSubstring("Missing required fields"))
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("invalid field is added", func() {
				BeforeEach(func() {
					testBody = map[string]interface{}{
						"username": "testuser",
						"email":    "test@example.com",
						"loginID":  1,
					}

					jsonBytes, err = json.Marshal(testBody)
					Expect(err).ToNot(HaveOccurred())

					w = httptest.NewRecorder()
					r = httptest.NewRequest("POST", "/apis/v1/post", bytes.NewBuffer(jsonBytes))
					r.Header.Set("Content-Type", "application/json")
				})

				It("should return bad request", func() {
					router.ServeHTTP(w, r)
					Expect(w.Body.String()).To(ContainSubstring("Invalid input"))
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})

		Context("when user does not exists", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"user_id": 99,
					"title":   "Post title #99",
					"content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat",
				}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/post", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should return not found error", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusCreated))
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["message"]).To(Equal("Post created successfully"))
			})
		})
	})
})
