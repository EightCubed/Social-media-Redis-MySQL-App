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

var _ = Describe("UserDelete", func() {
	var (
		testBody  map[string]interface{}
		router    *mux.Router
		w         *httptest.ResponseRecorder
		r         *http.Request
		jsonBytes []byte
		err       error
	)

	_, router = createFakeSocialMediaHandlerAndRouter()

	Context("when delete function is called", func() {
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

		Context("when user exists", func() {
			BeforeEach(func() {
				w = httptest.NewRecorder()
				r = httptest.NewRequest("DELETE", "/apis/v1/user/1", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})
			It("should return 200 OK and delete user", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when user does not exist", func() {
			BeforeEach(func() {
				w = httptest.NewRecorder()
				r = httptest.NewRequest("DELETE", "/apis/v1/user/11", bytes.NewBuffer(jsonBytes))
				r.Header.Set("Content-Type", "application/json")
			})
			It("should return 404 NOT FOUND", func() {
				router.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
