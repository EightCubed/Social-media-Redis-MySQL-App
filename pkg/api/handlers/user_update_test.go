package handlers_test

import (
	"bytes"
	"encoding/json"
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
		JustBeforeEach(func() {
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

			// Happy path
			Context("when request body is correct", func() {
				BeforeEach(func() {
					testBody = map[string]interface{}{
						"username": "testuser_updated",
						"email":    "test_updated@example.com",
						"password": "testpassword_updated",
					}

					jsonBytes, err = json.Marshal(testBody)
					Expect(err).ToNot(HaveOccurred())

					w = httptest.NewRecorder()
					r = httptest.NewRequest("PATCH", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
					r.Header.Set("Content-Type", "application/json")
				})

				It("should handle valid JSON body and return success", func() {
					router.ServeHTTP(w, r)
					Expect(w.Code).To(Equal(http.StatusOK))
					var responseBody map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &responseBody)
					Expect(err).ToNot(HaveOccurred())
					Expect(responseBody["Username"]).To(Equal("testuser"))
					redisResult, err := fakeSocialMediaHandler.RedisReader.Get("user:1").Result()
					var userObj userReturn
					Expect(err).ToNot(HaveOccurred())
					err = json.Unmarshal([]byte(redisResult), &userObj)
					Expect(err).ToNot(HaveOccurred())
					Expect(redisResult).ToNot(BeNil())
					Expect(userObj.Username).To(Equal("testuser_updated"))
					Expect(userObj.Email).To(Equal("test_updated@example.com"))
					Expect(userObj.ID).To(Equal(1))
					Expect(userObj.LoginID).To(Equal(1))
				})
			})

			Context("when user does not exist", func() {
				BeforeEach(func() {
					testBody = map[string]interface{}{
						"username": "testuser",
						"email":    "test@example.com",
						"password": "testpassword",
					}

					jsonBytes, err = json.Marshal(testBody)
					Expect(err).ToNot(HaveOccurred())

					w = httptest.NewRecorder()
					r = httptest.NewRequest("PATCH", "/apis/v1/user/12", bytes.NewBuffer(jsonBytes))
					r.Header.Set("Content-Type", "application/json")
				})

				It("should return not found error", func() {
					router.ServeHTTP(w, r)
					Expect(w.Code).To(Equal(http.StatusNotFound))
				})
			})

			Context("when request body is not correct", func() {
				BeforeEach(func() {
					w = httptest.NewRecorder()
					r = httptest.NewRequest("PATCH", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
					r.Header.Set("Content-Type", "application/json")
				})

				It("should return not bad request", func() {
					router.ServeHTTP(w, r)
					Expect(w.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})
	})
})
