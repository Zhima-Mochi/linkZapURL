package redirection

import (
	"context"
	"testing"
	"time"

	"github.com/Zhima-Mochi/linkZapURL/models"
	"github.com/Zhima-Mochi/linkZapURL/pkg/cache"
	"github.com/Zhima-Mochi/linkZapURL/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockCTX = context.Background()
)

type Mocks struct {
	mockDB    *database.MockDatabase
	mockCache *cache.MockCache
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		name  string
		code  string
		setUp func(mocks *Mocks)
		check func(t *testing.T, res *models.URL, err error)
	}{
		{
			name: "success",
			code: "5abcDEF",
			setUp: func(mocks *Mocks) {
				timeNow = func() time.Time {
					return time.Unix(10, 0)
				}
				mocks.mockCache.EXPECT().Get(mockCTX, "5abcDEF").Return(&models.URL{
					URL:      "https://www.google.com",
					ExpireAt: 1000,
				}, nil).Times(1)
			},
			check: func(t *testing.T, res *models.URL, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://www.google.com", res.URL)
				assert.Equal(t, int64(1000), res.ExpireAt)
			},
		},
		{
			name: "expired",
			code: "5abcDEF",
			setUp: func(mocks *Mocks) {
				timeNow = func() time.Time {
					return time.Unix(1001, 0)
				}
				mocks.mockCache.EXPECT().Get(mockCTX, "5abcDEF").Return(&models.URL{
					URL:      "https://www.google.com",
					ExpireAt: 1000,
				}, nil).Times(1)
			},
			check: func(t *testing.T, res *models.URL, err error) {
				assert.ErrorIs(t, err, ErrExpired)
			},
		},
		{
			name: "not found",
			code: "5abcDEF",
			setUp: func(mocks *Mocks) {
				timeNow = func() time.Time {
					return time.Unix(10, 0)
				}
				mocks.mockCache.EXPECT().Get(mockCTX, "5abcDEF").Return(nil, cache.ErrNotFound).Times(1)
				mocks.mockCache.EXPECT().Set(mockCTX, "5abcDEF", nil, nonExistedCacheTTL).Return(nil).Times(1)

				u := &models.URL{
					Code: "5abcDEF",
				}
				err := u.ToBSON()
				assert.NoError(t, err)

				mocks.mockDB.EXPECT().Get(mockCTX, collectionName, u.ID, u).Return(database.ErrNotFound).Times(1)
			},
			check: func(t *testing.T, res *models.URL, err error) {
				assert.ErrorIs(t, err, ErrNotFound)
			},
		},
		{
			name: "invalid code",
			code: "05abcDE",
			check: func(t *testing.T, res *models.URL, err error) {
				assert.ErrorIs(t, err, ErrInvalidCode)
			},
		},
	}

	for _, tt := range tests {
		originalTimeNow := timeNow
		defer func() { timeNow = originalTimeNow }()

		t.Run(tt.name, func(t *testing.T) {
			// Create a controller to manage the mock.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create the mocks.
			mocks := &Mocks{
				mockDB:    database.NewMockDatabase(ctrl),
				mockCache: cache.NewMockCache(ctrl),
			}

			// Create the object we are testing.
			redirection := NewRedirection(mocks.mockCache, mocks.mockDB)

			// Set up the mock expectations.
			if tt.setUp != nil {
				tt.setUp(mocks)
			}

			// Call the code we are testing.
			res, err := redirection.Redirect(mockCTX, tt.code)

			// Check the results.
			if tt.check != nil {
				tt.check(t, res, err)
			}
		})
	}
}
