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

	"go-social-media/pkg/api/handlers"
)

var _ = Describe("UserUpdate", func() {
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

	Context("when update function is called", func() {
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

		Context("when request body is valid", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username": "testuser_updated",
				}
				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("PATCH", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should return 200 OK and update user in Redis", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusOK))

				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				Expect(err).ToNot(HaveOccurred())
				Expect(responseBody["user_id"]).ToNot(BeNil())
				Expect(responseBody["message"]).To(ContainSubstring("User updated successfully"))

				redisResult, err := fakeSocialMediaHandler.RedisReader.Get("user:1").Result()
				Expect(err).ToNot(HaveOccurred())
				Expect(redisResult).ToNot(BeNil())

				fmt.Println(redisResult)

				var userObj userReturn
				err = json.Unmarshal([]byte(redisResult), &userObj)
				Expect(err).ToNot(HaveOccurred())
				Expect(userObj.Username).To(Equal("testuser_updated"))
				Expect(userObj.ID).To(Equal(1))
				Expect(userObj.LoginID).To(Equal(1))
			})
		})

		Context("when user does not exist", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username": "nonexistent",
				}
				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("PATCH", "/apis/v1/user/999", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should return 404 Not Found", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when request body has unexpected fields", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username":    "testuser",
					"email":       "test@example.com",
					"password":    "testpassword",
					"randomField": "should not be here",
				}
				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("PATCH", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should return 400 Bad Request", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when request body includes forbidden fields like loginID", func() {
			BeforeEach(func() {
				testBody = map[string]interface{}{
					"username": "testuser",
					"email":    "test@example.com",
					"password": "testpassword",
					"loginID":  1,
				}
				jsonBytes, err = json.Marshal(testBody)
				Expect(err).ToNot(HaveOccurred())

				w = httptest.NewRecorder()
				r = httptest.NewRequest("PATCH", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("should return 400 Bad Request", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
