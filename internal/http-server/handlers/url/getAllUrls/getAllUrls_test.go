package getAllUrls

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUrlGetter struct {
	mock.Mock
}

func (m *MockUrlGetter) GetAllUrls() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func TestGetAllUrlsHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUrlGetter)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success with urls",
			setupMock: func(m *MockUrlGetter) {
				m.On("GetAllUrls").Return(map[string]string{
					"abc": "https://example.com",
					"def": "https://google.com",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"result": {"abc":"https://example.com", "def":"https://google.com"},"status":"OK"}`,
		},
		{
			name: "empty result",
			setupMock: func(m *MockUrlGetter) {
				m.On("GetAllUrls").Return(map[string]string{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"OK","result":{}}`,
		},
		{
			name: "internal server error",
			setupMock: func(m *MockUrlGetter) {
				m.On("GetAllUrls").Return(map[string]string{}, errors.New("database error"))
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"error":"internal server error", "status":"Error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Создаем новый мок для каждого теста
			mockGetter := new(MockUrlGetter)
			tt.setupMock(mockGetter)

			log := slog.Default()

			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Get("/", New(log, mockGetter))

			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockGetter.AssertExpectations(t)
		})
	}
}
