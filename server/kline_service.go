package main

import (
	"time"

	"github.com/patrickmn/go-cache"

	"gin-example/model"
)

type KLineService struct {
	repo klineRepository

	cache *cache.Cache
}

func NewKLineSerice() *KLineService {
	c := cache.New(5*time.Minute, 10*time.Minute)

	return &KLineService{
		cache: c,
	}
}

func (s *KLineService) Get() (*model.KLine, error) {
	// 1. from cache
	// 2. from db
	return nil, nil
}
