package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zhima-Mochi/linkZapURL/handlers/redirection"
	"github.com/Zhima-Mochi/linkZapURL/handlers/shortening"
	"github.com/Zhima-Mochi/linkZapURL/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type Mocks struct {
	mockShortening  *shortening.MockShortening
	mockRedirection *redirection.MockRedirection
}

func NewMocks(ctrl *gomock.Controller) *Mocks {
	return &Mocks{
		mockShortening:  shortening.NewMockShortening(ctrl),
		mockRedirection: redirection.NewMockRedirection(ctrl),
	}
}

func TestShorten(t *testing.T) {
	originalEndpoint := endpoint
	defer func() {
		endpoint = originalEndpoint
	}()

	endpoint = "https://localhost"

	tests := []struct {
		name    string
		req     *ShortenRequest
		setUp   func(mocks *Mocks)
		check   func(t *testing.T, res *ShortenResponse)
		expCode int
	}{
		{
			name: "success",
			req: &ShortenRequest{
				URL:      "https://www.google.com",
				ExpireAt: "2021-02-08T09:20:41Z",
			},
			setUp: func(mocks *Mocks) {
				expireAtTime, _ := time.Parse(time.RFC3339, "2021-02-08T09:20:41Z")
				mocks.mockShortening.EXPECT().Shorten(gomock.Any(), "https://www.google.com", expireAtTime.Unix()).Return(&models.URL{
					ID:   1,
					Code: "1",
					URL:  "https://localhost/1",
				}, nil).Times(1)
			},
			check: func(t *testing.T, res *ShortenResponse) {
				assert.Equal(t, "1", res.ID)
				assert.Equal(t, "https://localhost/1", res.ShortURL)
			},
			expCode: http.StatusCreated,
		},
		{
			name: "empty url",
			req: &ShortenRequest{
				URL:      "",
				ExpireAt: "2021-02-08T09:20:41Z",
			},
			expCode: http.StatusBadRequest,
		},
		{
			name: "invalid expire at",
			req: &ShortenRequest{
				URL:      "https://www.google.com",
				ExpireAt: "1",
			},
			expCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mocks := NewMocks(ctrl)

		h := &Handler{
			shortening:  mocks.mockShortening,
			redirection: mocks.mockRedirection,
		}

		router.POST("/api/v1/urls", h.Shorten)

		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp(mocks)
			}

			reqBody, err := json.Marshal(tt.req)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expCode, w.Code)

			if tt.check != nil {
				res := &ShortenResponse{}
				err := json.Unmarshal(w.Body.Bytes(), res)
				assert.NoError(t, err)
				tt.check(t, res)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	originalEndpoint := endpoint
	defer func() {
		endpoint = originalEndpoint
	}()

	endpoint = "https://localhost"

	tests := []struct {
		name    string
		code    string
		setUp   func(mocks *Mocks)
		expCode int
	}{
		{
			name: "success",
			code: "1",
			setUp: func(mocks *Mocks) {
				mocks.mockRedirection.EXPECT().Redirect(gomock.Any(), "1").Return(&models.URL{
					ID:   1,
					Code: "1",
					URL:  "https://www.google.com",
				}, nil).Times(1)
			},
			expCode: http.StatusMovedPermanently,
		},
		{
			name: "not found",
			code: "1",
			setUp: func(mocks *Mocks) {
				mocks.mockRedirection.EXPECT().Redirect(gomock.Any(), "1").Return(nil, redirection.ErrNotFound).Times(1)
			},
			expCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		gin.SetMode(gin.TestMode)
		router := gin.Default()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mocks := NewMocks(ctrl)

		h := &Handler{
			shortening:  mocks.mockShortening,
			redirection: mocks.mockRedirection,
		}

		router.GET("/:code", h.Redirect)

		t.Run(tt.name, func(t *testing.T) {
			if tt.setUp != nil {
				tt.setUp(mocks)
			}

			req := httptest.NewRequest("GET", "/"+tt.code, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expCode, w.Code)
		})
	}
}
