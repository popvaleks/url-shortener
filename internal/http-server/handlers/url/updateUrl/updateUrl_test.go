package updateUrl

import (
	"errors"
	"github.com/popvaleks/url-shortener/internal/storage"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUrlEditer struct {
	mock.Mock
}

func (m *MockUrlEditer) UpdateUrl(url, alias string) (string, error) {
	args := m.Called(url, alias)
	return args.String(0), args.Error(1)
}

func TestUpdateUrlHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		alias          string
		setupMock      func(*MockUrlEditer)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success update",
			requestBody: `{"url": "http://example.com"}`,
			alias:       "test",
			setupMock: func(m *MockUrlEditer) {
				m.On("UpdateUrl", "http://example.com", "test").Return("test", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"result": {"alias":"test"},"status":"OK"}`,
		},
		{
			name:        "alias not found",
			requestBody: `{"url": "http://example.com"}`,
			alias:       "notFound",
			setupMock: func(m *MockUrlEditer) {
				m.On("UpdateUrl", "http://example.com", "notFound").Return("", storage.ErrAliasNotFound)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"error":"alias not found", "status":"Error"}`,
		},
		{
			name:        "internal server error",
			requestBody: `{"url": "http://example.com"}`,
			alias:       "err",
			setupMock: func(m *MockUrlEditer) {
				m.On("UpdateUrl", "http://example.com", "err").Return("", errors.New("internal error"))
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "Error","error": "failed to update url"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockEditer := new(MockUrlEditer)
			log := slog.Default()
			tt.setupMock(mockEditer)
			t.Parallel()

			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Patch("/{alias}", New(log, mockEditer))

			req, err := http.NewRequest("PATCH", "/"+tt.alias, strings.NewReader(tt.requestBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockEditer.AssertExpectations(t)
		})
	}
}
