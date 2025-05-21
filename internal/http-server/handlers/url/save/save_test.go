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
	log := slog.Default()

	tests := []struct {
		name         string
		requestBody  string
		setupMock    func(*MockUrlSaver) // Принимает конкретный мок
		expectedCode int
		expectedBody string
	}{
		{
			name:        "success with alias",
			requestBody: `{"url": "http://example.com", "alias": "example"}`,
			setupMock: func(m *MockUrlSaver) {
				m.On("SaveUrl", "http://example.com", "example").Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"OK","alias":"example"}`,
		},
		{
			name:        "success without alias",
			requestBody: `{"url": "http://example.com"}`,
			setupMock: func(m *MockUrlSaver) {
				m.On("SaveUrl", "http://example.com", mock.AnythingOfType("string")).Return(int64(1), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"OK","alias":`,
		},
		{
			name:        "url already exists",
			requestBody: `{"url": "http://example.com", "alias": "example"}`,
			setupMock: func(m *MockUrlSaver) {
				m.On("SaveUrl", "http://example.com", "example").Return(int64(0), storage.ErrUrlExists)
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"Error","error":"url already exists"}`,
		},
		{
			name:         "invalid url",
			requestBody:  `{"url": "invalid-url"}`,
			setupMock:    func(m *MockUrlSaver) {},
			expectedCode: http.StatusOK,
			expectedBody: `{"status":"Error","error":"field Url is not a valid URL"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Создаем новый мок для каждого теста
			mockSaver := new(MockUrlSaver)
			tt.setupMock(mockSaver) // Передаем конкретный мок

			req, err := http.NewRequest("POST", "/url", strings.NewReader(tt.requestBody))
			assert.NoError(t, err)

			handler := middleware.RequestID(http.HandlerFunc(New(log, mockSaver)))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}

			// Проверяем ожидания только для текущего мока
			mockSaver.AssertExpectations(t)
		})
	}
}
