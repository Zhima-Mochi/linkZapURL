package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	workDir, err := os.Getwd()
	assert.NoError(t, err)
	fileName := "config.yaml"

	tests := []struct {
		name   string
		data   []byte
		expRes *Config
	}{
		{
			name: "mongodb",
			data: []byte(`
mongodb:
  uri: "mongodb://localhost:27017"
`),
			expRes: &Config{
				Mongodb: &Mongodb{
					URI: "mongodb://localhost:27017",
				},
			},
		},
		{
			name: "redis",
			data: []byte(`
redis:
  addrs:
    - "localhost:6379"
  password: "password"
  db: 0
`),
			expRes: &Config{
				Redis: &Redis{
					Addrs:    []string{"localhost:6379"},
					Password: "password",
					DB:       0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp(workDir, fileName)
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.Write(tt.data)
			assert.NoError(t, err)

			configFile = tmpFile.Name()

			config, err := GetConfig()
			assert.NoError(t, err)
			assert.Equal(t, tt.expRes, config)
		})
	}
}
