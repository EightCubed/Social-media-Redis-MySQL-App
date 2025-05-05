package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-social-media/pkg/api/handlers"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type UserObject struct {
	ID        int    `json:"ID"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
	DeletedAt string `json:"DeletedAt"`
	Username  string `json:"Username"`
	Email     string `json:"Email"`
	LoginID   int    `json:"LoginID"`
}

var _ = Describe("UserPost", func() {
	var (
		testBody               map[string]interface{}
		fakeSocialMediaHandler *handlers.SocialMediaHandler
		router                 *mux.Router
		w                      *httptest.ResponseRecorder
		r                      *http.Request
		jsonBytes              []byte
		err                    error
	)

	fakeSocialMediaHandler, router = createFakeSocialMediaHandlerAndRouter()

	// Happy path
	Context("when request body is valid", func() {
		BeforeEach(func() {
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
		})

		It("should handle valid JSON body and return success", func() {
			router.ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusCreated))
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseBody["message"]).To(Equal("User created successfully"))
			redisResult, err := fakeSocialMediaHandler.RedisReader.Get("user:1").Result()
			var userObj UserObject
			Expect(err).ToNot(HaveOccurred())
			err = json.Unmarshal([]byte(redisResult), &userObj)
			Expect(err).ToNot(HaveOccurred())
			Expect(redisResult).ToNot(BeNil())
			Expect(userObj.Username).To(Equal("testuser"))
			Expect(userObj.Email).To(Equal("test@example.com"))
			Expect(userObj.ID).To(Equal(1))
			Expect(userObj.LoginID).To(Equal(1))
		})
	})

	Context("when request body is not correct", func() {
		Context("request body is nil", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should return Bad request", func() {
				router.ServeHTTP(w, r)
				Expect(w.Body.String()).To(ContainSubstring("Missing required fields"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
		Context("username is nil", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"email":    "test@example.com",
					"password": "testpassword",
				}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should handle valid JSON body and return success", func() {
				router.ServeHTTP(w, r)
				Expect(w.Body.String()).To(ContainSubstring("Missing required fields"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
		Context("email is nil", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username": "testuser",
					"password": "testpassword",
				}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should handle valid JSON body and return success", func() {
				router.ServeHTTP(w, r)
				Expect(w.Body.String()).To(ContainSubstring("Missing required fields"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
		Context("password is nil", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
				}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should handle valid JSON body and return success", func() {
				router.ServeHTTP(w, r)
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
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should handle valid JSON body and return success", func() {
				router.ServeHTTP(w, r)
				Expect(w.Body.String()).To(ContainSubstring("Missing required fields"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("duplicate user is added field is added", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"loginID":  1,
				}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")

				testBody = map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"loginID":  1,
				}

				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("POST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should handle return error", func() {
				router.ServeHTTP(w, r)
				Expect(w.Body.String()).To(ContainSubstring("Missing required fields"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
