package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-social-media/pkg/api/handlers"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserPost", func() {
	var (
		testBody               map[string]interface{}
		fakeSocialMediaHandler *handlers.SocialMediaHandler
		w                      *httptest.ResponseRecorder
		r                      *http.Request
		jsonBytes              []byte
		err                    error
	)

	BeforeEach(func() {
		fakeSocialMediaHandler = createFakeSocialMediaHandler()
		w = httptest.NewRecorder()
	})

	// Happy path
	Context("when body is valid", func() {
		BeforeEach(func() {
			testBody = map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			}

			jsonBytes, err = json.Marshal(testBody)
			Expect(err).ToNot(HaveOccurred())

			r = httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBytes))
			r.Header.Set("Content-Type", "application/json")
		})

		It("should handle valid JSON body and return success", func() {
			fakeSocialMediaHandler.PostUser(w, r)

			Expect(w.Code).To(Equal(http.StatusCreated))

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			Expect(err).ToNot(HaveOccurred())

			Expect(responseBody["success"]).To(BeTrue())
			Expect(responseBody["username"]).To(Equal("testuser"))
		})
	})
})
