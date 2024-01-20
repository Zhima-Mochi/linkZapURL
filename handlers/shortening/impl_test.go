package shortening

import (
	"context"
	"testing"

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
		name   string
		url    string
		setUp  func(mocks *Mocks)
		expErr error
	}{
		{
			name: "success",
			url:  "https://www.google.com",
			setUp: func(mocks *Mocks) {
				mocks.mockDB.EXPECT().Set(mockCTX, collectionName, gomock.Any(), "https://www.google.com").Return(nil).Times(1)
			},
		},
		{
			name:   "empty url",
			url:    "",
			expErr: ErrEmptyURL,
		},
		{
			name: "collision",
			url:  "https://www.google.com",
			setUp: func(mocks *Mocks) {
				mocks.mockDB.EXPECT().Set(mockCTX, collectionName, gomock.Any(), "https://www.google.com").Return(database.ErrCollision).Times(1)
			},
			expErr: database.ErrCollision,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mocks := &Mocks{
				mockDB: database.NewMockDatabase(ctrl),
			}
			shortening := NewShortening(mockMachineID, mocks.mockDB)

			if tt.setUp != nil {
				tt.setUp(mocks)
			}

			_, err := shortening.Shorten(mockCTX, tt.url)
			assert.ErrorIs(t, err, tt.expErr)
		})
	}
}
