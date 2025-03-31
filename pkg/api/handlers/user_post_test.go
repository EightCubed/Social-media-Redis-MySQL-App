package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-social-media/pkg/api/handlers"
	"go-social-media/pkg/models"
	"go-social-media/pkg/utils"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// MockDBWriter is a mock implementation of the database writer
type MockDBWriter struct {
	mock.Mock
}

func (m *MockDBWriter) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDBWriter) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// MockRedisReader is a mock implementation of the Redis cache
type MockRedisReader struct {
	mock.Mock
}

func (m *MockRedisReader) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func TestPostUser(t *testing.T) {
	testCases := []struct {
		name               string
		inputBody          interface{}
		mockDBCreateLogin  func(*MockDBWriter)
		mockDBCreateUser   func(*MockDBWriter)
		mockRedisSet       func(*MockRedisReader)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Successful User Creation",
			inputBody: handlers.UserPostBody{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockDBCreateLogin: func(mockDB *MockDBWriter) {
				mockDB.On("Create", mock.AnythingOfType("*models.Login")).
					Return(&gorm.DB{}).
					Run(func(args mock.Arguments) {
						login := args.Get(0).(*models.Login)
						login.ID = 1 // Simulate DB assigning an ID
					})
			},
			mockDBCreateUser: func(mockDB *MockDBWriter) {
				mockDB.On("Create", mock.AnythingOfType("*models.User")).
					Return(&gorm.DB{}).
					Run(func(args mock.Arguments) {
						user := args.Get(0).(*models.User)
						user.ID = 1 // Simulate DB assigning an ID
					})
			},
			mockRedisSet: func(mockRedis *MockRedisReader) {
				mockRedis.On("Set", mock.AnythingOfType("string"), mock.Anything, mock.Anything).
					Return(&redis.StatusCmd{})
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"message":"User created successfully"}`,
		},
		{
			name: "Invalid Input",
			inputBody: map[string]interface{}{
				"invalid": "input",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Invalid input\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare mock objects
			mockDBWriter := new(MockDBWriter)
			mockRedisReader := new(MockRedisReader)

			// Set up mock expectations if provided
			if tc.mockDBCreateLogin != nil {
				tc.mockDBCreateLogin(mockDBWriter)
			}
			if tc.mockDBCreateUser != nil {
				tc.mockDBCreateUser(mockDBWriter)
			}
			if tc.mockRedisSet != nil {
				tc.mockRedisSet(mockRedisReader)
			}

			// Create handler with mock dependencies
			handler := &handlers.SocialMediaHandler{
				DBWriter:    mockDBWriter,
				RedisReader: mockRedisReader,
			}

			// Convert input body to JSON
			jsonBody, err := json.Marshal(tc.inputBody)
			assert.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			w := httptest.NewRecorder()

			// Call the handler
			handler.PostUser(w, req)

			// Check response status code
			assert.Equal(t, tc.expectedStatusCode, w.Code)

			// Check response body
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())

			// Verify mock expectations
			mockDBWriter.AssertExpectations(t)
			mockRedisReader.AssertExpectations(t)
		})
	}
}

// Additional test cases to consider
func TestPasswordHashing(t *testing.T) {
	inputPassword := "testpassword"

	// Simulate user creation
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// Verify password can be compared
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(inputPassword))
	assert.NoError(t, err)
}

func TestTokenGeneration(t *testing.T) {
	userID := uint(1)

	accessToken, err := utils.GenerateAccessToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	refreshToken, err := utils.GenerateRefreshToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
}
