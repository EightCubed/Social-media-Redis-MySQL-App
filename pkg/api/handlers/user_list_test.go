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

var _ = Describe("UserList", func() {
	var (
		testBody map[string]interface{}
		// fakeSocialMediaHandler *handlers.SocialMediaHandler
		router    *mux.Router
		w         *httptest.ResponseRecorder
		r         *http.Request
		jsonBytes []byte
		err       error
	)

	_, router = createFakeSocialMediaHandlerAndRouter()

	Context("when list function is called", func() {
		Context("when list is empty", func() {
			BeforeEach(func() {
				w = httptest.NewRecorder()
				r = httptest.NewRequest("LIST", "/apis/v1/user", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})

			It("return expected output", func() {
				router.ServeHTTP(w, r)
			})
		})

		Context("when list is not empty", func() {
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
			})

			It("return expected output", func() {
				router.ServeHTTP(w, r)
			})
		})
	})
})
