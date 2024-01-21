package redirection

import (
	"context"
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
	now := timeNow().Unix()

	// Check if the code is in the cache.
	if b, err := im.cache.Get(ctx, code); err == nil {
		val := b.(*models.URL)

		if val.ExpireAt < now {
			return nil, ErrExpired
		}

		return val, nil
	}

	// Check if the code is in the database.
	doc := &models.URL{
		Code: code,
	}

	id, err := doc.FillID()
	if err != nil {
		return nil, err
	}

	err = im.database.Get(ctx, collectionName, id, doc)
	if err == database.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Set the code in the cache.
	err = im.cache.Set(ctx, code, doc, cacheTTL)
	if err != nil {
		return nil, err
	}

	if doc.ExpireAt < now {
		return nil, ErrExpired
	}

	return doc, nil
}
