package services

import (
	"caching"
	"context"
)

type CacheService struct {
	Cache    caching.Cache
	Context  context.Context
	Language string
}
