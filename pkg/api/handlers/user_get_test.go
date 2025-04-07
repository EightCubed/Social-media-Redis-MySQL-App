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

var _ = Describe("UserGet", func() {
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

	Context("when get function is called", func() {
		Context("when user exists", func() {
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
				router.ServeHTTP(w, r)

				w = httptest.NewRecorder()
				r = httptest.NewRequest("LIST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(w, r)

				w = httptest.NewRecorder()
				r = httptest.NewRequest("GET", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("returns a successful response", func() {
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
			})
		})
		Context("when user does not exist", func() {
			BeforeEach(func() {
				w = httptest.NewRecorder()
				r = httptest.NewRequest("GET", "/apis/v1/user/12", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("returns a not found response", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
