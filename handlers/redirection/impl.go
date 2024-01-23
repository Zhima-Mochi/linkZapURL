package redirection

import (
	"context"
	"log"
	"time"

	"github.com/Zhima-Mochi/linkZapURL/models"
	"github.com/Zhima-Mochi/linkZapURL/pkg/cache"
	"github.com/Zhima-Mochi/linkZapURL/pkg/database"
)

const (
	collectionName = "url"
)

var (
	timeNow = time.Now

	cacheTTL = 1 * time.Hour

	nonExistedCacheTTL = 5 * time.Minute
)

type impl struct {
	cache    cache.Cache
	database database.Database
}

func NewRedirection(cache cache.Cache, database database.Database) Redirection {
	return &impl{
		cache:    cache,
		database: database,
	}
}

// Redirect redirects the shortened url to the original url.
func (im *impl) Redirect(ctx context.Context, code string) (*models.URL, error) {
	// Validate the code.
	if _, err := models.Decode(code); err != nil {
		return nil, ErrInvalidCode
	}

	now := timeNow().Unix()

	// Check if the code is in the cache.
	u := &models.URL{}
	if err := im.cache.Get(ctx, code, u); err == nil {
		log.Println("Cache hit")
		// non-existent codes are set to nil in the cache.
		if u.Code == "" {
			return nil, ErrNotFound
		}

		if u.ExpireAt < now {
			return nil, ErrExpired
		}

		return u, nil
	}

	log.Println("Cache miss")

	// Check if the code is in the database.
	u = &models.URL{
		Code: code,
	}

	err := u.ToModel()
	if err != nil {
		return nil, err
	}

	err = im.database.Get(ctx, collectionName, u.ID, u)
	if err == database.ErrNotFound {
		// Set non-existent codes in the cache to prevent hitting the database again.
		err := im.cache.Set(ctx, code, nil, nonExistedCacheTTL)
		if err != nil {
			log.Println("Failed to set non-existent code in cache:", err)
		}

		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Set the code in the cache.
	err = im.cache.Set(ctx, code, u, cacheTTL)
	if err != nil {
		return nil, err
	}

	if u.ExpireAt < now {
		return nil, ErrExpired
	}

	err = u.ToJSON()
	if err != nil {
		return nil, err
	}

	return u, nil
}
