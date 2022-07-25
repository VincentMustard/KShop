package main

import (
	"database/sql"
	"time"

	"gin-example/model"
)

type klineRepository interface {
	Get(startTime time.Time, interval time.Duration, kind string) (*model.KLine, error)
	Create(kline *model.KLine) error
}

type klineRepositoryImpl struct {
	db *sql.DB
}

func NewKlineRepository() klineRepository {
	return &klineRepositoryImpl{}
}

func (r *klineRepositoryImpl) Create(kline *model.KLine) error {
	return nil
}

func (r *klineRepositoryImpl) Get(
	startTime time.Time, interval time.Duration, kind string,
) (*model.KLine, error) {
	return nil, nil
}
