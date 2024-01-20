package shortening

import (
	"context"
	"sync"
	"time"

	"github.com/Zhima-Mochi/linkZapURL/pkg/database"
)

const (
	base58alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

	collectionName = "shortening"
)

var (
	timeNow = time.Now

	seq uint8

	encode = func(num int64) string {
		codes := make([]byte, 0, 7)

		for num > 0 {
			codes = append(codes, base58alphabet[num%58])
			num /= 58
		}

		return string(codes)
	}
)

type impl struct {
	machineID     uint8
	lastTimestamp int64
	database      database.Database
	mutex         sync.Mutex
}

func NewShortening(machineID uint8, database database.Database) Shortening {
	return &impl{
		machineID: machineID,
		database:  database,
	}
}

// nextSeq generates the next sequence number.
func (im *impl) nextSeq(now int64) uint8 {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if now == im.lastTimestamp {
		seq++
	} else {
		seq = 0
		im.lastTimestamp = now
	}

	return seq
}

// generateID generates a unique ID based on the current time, machine ID and sequence number.
func (im *impl) generateID() string {
	now := timeNow().UnixMilli()
	num := int64(0)
	num |= (now & 0x1FFFFFF) << 16
	num |= int64(im.machineID) << 8
	num |= int64(im.nextSeq(now))

	return encode(num)
}

func (im *impl) Shorten(ctx context.Context, url string) (string, error) {
	if url == "" {
		return "", ErrEmptyURL
	}

	id := im.generateID()

	err := im.database.Set(ctx, collectionName, id, url)
	if err != nil {
		return "", err
	}

	return id, nil
}
