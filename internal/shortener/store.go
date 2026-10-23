package shortener

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("URL mapping not found")

type Store interface {
	Save(ctx context.Context,
		mapping URLMapping) error
	Get(ctx context.Context, code string) (URLMapping, error)
}
