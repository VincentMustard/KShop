package main

import (
	"context"
	"gin-example/model"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Consumer struct {
	rdb *redis.Client

	orders []*model.Order
	mu     sync.Mutex
	timer  *time.Ticker
	repo   klineRepository
	logger *zap.Logger

	latestK1   *model.KLine
	latestK1Mu sync.Mutex

	latestK5   *model.KLine
	latestK5Mu sync.Mutex

	klinePublishTimer *time.Ticker
}

func NewComsumer(Addr string, Password string, DB int) *Consumer {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password, // "" means no password set
		DB:       DB,       // 0 means using default DB
	})

	logger, _ := zap.NewProduction()

	return &Consumer{
		rdb:               rdb,
		logger:            logger,
		klinePublishTimer: time.NewTicker(time.Millisecond * 500), //将最新k线发布的频率设置为0.5s
	}
}

func (c *Consumer) Consume(SubChannel string) {
	pubsub := c.rdb.Subscribe(context.TODO(), SubChannel)
	defer pubsub.Close()

	// 创建goroutine，使其通过consumer中的交易对创建K线并存入mysql
	go c.GenerateKLine()

	// 周期性发布最新k线到redis中
	go c.publishK1linePeriodically("latestK1")
	go c.publishK5linePeriodically("latestK5")

	// 从Redis中订阅交易对并存入consumer中
	for {
		msg, err := pubsub.ReceiveMessage(context.TODO())
		if err != nil {
			c.logger.Sugar().Errorf("failed to receive message: %v", err)
		}

		o, err := getOrderFromPayload(msg.Payload)
		if err != nil {
			c.logger.Sugar().Errorf("failed to get order from payload: %v", err)
		}

		c.AddOrder(o)
	}
}

func (c *Consumer) GenerateKLine() {
	for {
		select {
		case <-c.timer.C:
			now := time.Now().UTC()

			if now.Second() == 1 { //每分钟的第一秒生成一次k1
				orders := c.getOrders(now, time.Minute*1)
				k1 := generateKline(orders) //生成k线
				err := c.repo.Create(k1)    //将k线存储到mysql中
				if err != nil {
					c.logger.Sugar().Errorf("failed to create kline: %v", err)
				}

				//将k线发布到redis中
				c.mu.Lock()
				_ = c.rdb.Publish(context.TODO(), "NewKline", k1)
				c.mu.Unlock()

				c.latestK1Mu.Lock()
				c.latestK1 = k1
				c.latestK1Mu.Unlock()

				if now.Minute()%5 == 0 { //每5分钟的第一秒生成一次k5
					orders := c.getOrders(now, time.Minute*5)
					k5 := generateKline(orders)
					err := c.repo.Create(k5)
					if err != nil {
						c.logger.Sugar().Errorf("failed to create kline: %v", err)
					}

					//将k线发布到redis中
					c.mu.Lock()
					_ = c.rdb.Publish(context.TODO(), "NewKline", k5)
					c.mu.Unlock()

					c.latestK5Mu.Lock()
					c.latestK5 = k5
					c.latestK5Mu.Unlock()
				}
			}
		}
	}
}

func (c *Consumer) getOrders(now time.Time, interval time.Duration) []*model.Order {
	var slice, newOrders []*model.Order

	if interval == time.Minute*1 {
		for _, v := range c.orders {
			if (now.Minute() == 0 && v.Time.Minute() == 59) || v.Time.Minute() == now.Minute()-1 {
				slice = append(slice, v)
			}
		}
	}

	if interval == time.Minute*5 {
		for _, v := range c.orders {
			if (now.Minute() == 0 && v.Time.Minute() >= 55) || v.Time.Minute()-now.Minute() <= 5 {
				slice = append(slice, v)
			} else {
				newOrders = append(newOrders, v)
			}
		}

		c.orders = newOrders
	}

	return slice
}

func generateKline(orders []*model.Order) *model.KLine {
	var Open, Close, High, Low float64
	if len(orders) == 1 {
		Open = orders[0].Price
		Close = orders[0].Price
		High = orders[0].Price
		Low = orders[0].Price
	} else {
		Open = orders[0].Price
		High = orders[0].Price
		Low = orders[0].Price
		Close = orders[len(orders)-1].Price

		for _, v := range orders {
			price := v.Price
			if price > High {
				High = price
			}
			if price < Low {
				Low = price
			}
		}
	}

	return &model.KLine{
		Open:  Open,
		Close: Close,
		High:  High,
		Low:   Low,
	}
}

func getOrderFromPayload(string) (*model.Order, error) {
	return nil, nil
}

func (c *Consumer) AddOrder(element *model.Order) {
	c.mu.Lock()
	c.orders = append(c.orders, element)
	c.mu.Unlock()
}

func (c *Consumer) publishK1linePeriodically(channel string) {
	for {
		select {
		case <-c.klinePublishTimer.C:
			if c.latestK1 != nil {
				c.publishK1line(channel)
			}
		}
	}

}

func (c *Consumer) publishK5linePeriodically(channel string) {
	for {
		select {
		case <-c.klinePublishTimer.C:
			if c.latestK5 != nil {
				c.publishK5line(channel)
			}
		}
	}

}

func (c *Consumer) publishK1line(channel string) {
	c.latestK1Mu.Lock()
	defer c.latestK1Mu.Unlock()

	_ = c.rdb.Publish(context.TODO(), channel, c.latestK1)
}

func (c *Consumer) publishK5line(channel string) {
	c.latestK5Mu.Lock()
	defer c.latestK5Mu.Unlock()

	_ = c.rdb.Publish(context.TODO(), channel, c.latestK5)
}
