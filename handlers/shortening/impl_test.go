package shortening

import (
	"context"
	"testing"
	"time"

	"github.com/Zhima-Mochi/linkZapURL/models"
	"github.com/Zhima-Mochi/linkZapURL/pkg/database"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	mockCTX = context.Background()

	mockMachineID = uint8(1)
)

type Mocks struct {
	mockDB *database.MockDatabase
}

func TestShorten(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expireAt int64
		setUp    func(mocks *Mocks)
		check    func(t *testing.T, res *models.URL, err error)
	}{
		{
			name:     "success",
			url:      "https://www.google.com",
			expireAt: 1000,
			setUp: func(mocks *Mocks) {
				timeNow = func() time.Time {
					return time.Unix(10, 0)
				}
				mocks.mockDB.EXPECT().Set(mockCTX, collectionName, gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			check: func(t *testing.T, res *models.URL, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://www.google.com", res.URL)
				assert.Equal(t, int64(1000), res.ExpireAt)
			},
		},
		{
			name: "empty url",
			url:  "",
			check: func(t *testing.T, res *models.URL, err error) {
				assert.ErrorIs(t, err, ErrEmptyURL)
			},
		},
		{
			name:     "invalid expire at",
			url:      "https://www.google.com",
			expireAt: 0,
			setUp: func(mocks *Mocks) {
				timeNow = func() time.Time {
					return time.Unix(10, 0)
				}
			},
			check: func(t *testing.T, res *models.URL, err error) {
				assert.ErrorIs(t, err, ErrInvalidExpireAt)
			},
		},
		{
			name:     "collision",
			url:      "https://www.google.com",
			expireAt: 1000,
			setUp: func(mocks *Mocks) {
				timeNow = func() time.Time {
					return time.Unix(10, 0)
				}
				mocks.mockDB.EXPECT().Set(mockCTX, collectionName, gomock.Any(), gomock.Any()).Return(database.ErrCollision).Times(1)
			},
			check: func(t *testing.T, res *models.URL, err error) {
				assert.ErrorIs(t, err, database.ErrCollision)
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
				mockDB: database.NewMockDatabase(ctrl),
			}

			// Create the object we are testing.
			shortening := NewShortening(mockMachineID, mocks.mockDB)

			// Set up the mock expectations.
			if tt.setUp != nil {
				tt.setUp(mocks)
			}

			// Call the code we are testing.
			res, err := shortening.Shorten(mockCTX, tt.url, tt.expireAt)

			// Check the results.
			if tt.check != nil {
				tt.check(t, res, err)
			}
		})
	}
}
