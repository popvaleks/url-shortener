package save

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/popvaleks/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUrlSaver struct {
	mock.Mock
}

func (m *MockUrlSaver) SaveUrl(inputUrl string, alias string) (int64, error) {
	args := m.Called(inputUrl, alias)
	return args.Get(0).(int64), args.Error(1)
}

func TestSaveHandler(t *testing.T) {
	mockSaver := new(MockUrlSaver)
	log := slog.Default()

	tests := []struct {
		name         string
		requestBody  string
		setupMock    func()
		expectedCode int
		expectedBody string
	}{
		{
			name:        "success with alias",
			requestBody: `{"url": "http://example.com", "alias": "example"}`,
			setupMock: func() {
				mockSaver.On("SaveUrl", "http://example.com", "example").Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"Ok","alias":"example"}`,
		},
		{
			name:        "success without alias",
			requestBody: `{"url": "http://example.com"}`,
			setupMock: func() {
				mockSaver.On("SaveUrl", "http://example.com", mock.AnythingOfType("string")).Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"Ok","alias":`,
		},
		{
			name:        "url already exists",
			requestBody: `{"url": "http://example.com", "alias": "example"}`,
			setupMock: func() {
				mockSaver.On("SaveUrl", "http://example.com", "example").Return(int64(0), storage.ErrUrlExists)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"Ok","alias":"example"}`,
		},
		{
			name:         "invalid url",
			requestBody:  `{"url": "invalid-url"}`,
			setupMock:    func() {},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"Error","error":"field Url is not a valid URL"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			t.Parallel()

			req, err := http.NewRequest("POST", "/url", strings.NewReader(tt.requestBody))
			assert.NoError(t, err)

			// Добавляем middleware.RequestID для добавления ID запроса в контекст
			handler := middleware.RequestID(http.HandlerFunc(New(log, mockSaver)))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}

			mockSaver.AssertExpectations(t)
		})
	}
}
