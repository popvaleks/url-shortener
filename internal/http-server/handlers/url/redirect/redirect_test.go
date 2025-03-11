package redirect

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

type MockUrlGetter struct {
	mock.Mock
}

func (m *MockUrlGetter) GetUrl(alias string) (string, error) {
	args := m.Called(alias)
	return args.String(0), args.Error(1)
}

func TestRedirectHandler(t *testing.T) {
	mockGetter := new(MockUrlGetter)
	log := slog.Default()

	tests := []struct {
		name           string
		alias          string
		setupMock      func()
		expectedStatus int
		expectedURL    string
	}{
		{
			name:  "success",
			alias: "example",
			setupMock: func() {
				mockGetter.On("GetUrl", "example").Return("http://example.com", nil)
			},
			expectedStatus: http.StatusFound,
			expectedURL:    "http://example.com",
		},
		{
			name:  "alias not allowed",
			alias: "",
			setupMock: func() {
			},
			expectedStatus: http.StatusNotFound,
			expectedURL:    "",
		},
		{
			name:  "url not found",
			alias: "notfound",
			setupMock: func() {
				mockGetter.On("GetUrl", "notfound").Return("", storage.ErrUrlNotFound)
			},
			expectedStatus: http.StatusOK,
			expectedURL:    "",
		},
		{
			name:  "internal server error",
			alias: "error",
			setupMock: func() {
				mockGetter.On("GetUrl", "error").Return("", errors.New("internal error"))
			},
			expectedStatus: http.StatusOK,
			expectedURL:    "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			t.Parallel()

			r := chi.NewRouter()
			r.Use(middleware.RequestID)
			r.Get("/{alias}", New(log, mockGetter))

			req, err := http.NewRequest("GET", "/"+tt.alias, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedURL != "" {
				assert.Equal(t, tt.expectedURL, rr.Header().Get("Location"))
			} else {
				// Проверяем, что редиректа не было
				assert.Empty(t, rr.Header().Get("Location"))
			}

			mockGetter.AssertExpectations(t)
		})
	}
}
