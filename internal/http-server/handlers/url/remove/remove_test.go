package remove

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/popvaleks/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUrlRemover struct {
	mock.Mock
}

func (m *MockUrlRemover) DeleteUrl(alias string) error {
	args := m.Called(alias)
	return args.Error(0)
}

func TestRemoveHandler(t *testing.T) {
	mockRemover := new(MockUrlRemover)
	log := slog.Default()

	tests := []struct {
		name           string
		alias          string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:  "success",
			alias: "example",
			setupMock: func() {
				mockRemover.On("DeleteUrl", "example").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"OK"}`,
		},
		{
			name:  "url not found",
			alias: "notfound",
			setupMock: func() {
				mockRemover.On("DeleteUrl", "notfound").Return(storage.ErrUrlNotFound)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"Error","error":"url not found"}`,
		},
		{
			name:  "internal server error",
			alias: "error",
			setupMock: func() {
				mockRemover.On("DeleteUrl", "error").Return(errors.New("internal error"))
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"Error","error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Delete("/{alias}", New(log, mockRemover))

			req, err := http.NewRequest("DELETE", "/"+tt.alias, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockRemover.AssertExpectations(t)
		})
	}
}
