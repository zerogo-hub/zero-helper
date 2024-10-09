package cache

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Publish 将信息 message 发送到指定的频道 channel
// return 接收到信息 message 的订阅者数量
func (c *cache) Publish(channel string, message interface{}) (int, error) {
	return c.Int(c.DO("PUBLISH", channel, message))
}

// Subscribe 订阅给定的一个或多个频道的信息
// onReady 所有频道都订阅成功时调用，可选
// onMessage 接收到信息时调用
// num1 订阅失败后重试次数
// num2 发生异常后重试次数，-1 表示一直重试
func (c *cache) Subscribe(onReady func() error, onMessage func(channel string, data []byte) error, num1, num2 int, channels ...string) error {
	if onMessage == nil {
		return errors.New("onMessage cant be nil")
	}

	psc := redis.PubSubConn{Conn: c.pool.Get()}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		if num1 > 0 {
			time.Sleep(time.Second)
			return c.Subscribe(onReady, onMessage, num1-1, num2, channels...)
		}
		return err
	}

	quit := make(chan error, 1)

	go func() {
		for {
			switch v := psc.Receive().(type) {
			case error:
				quit <- v
				return
			case redis.Message:
				if err := onMessage(v.Channel, v.Data); err != nil {
					quit <- err
					return
				}
			case redis.Subscription:
				switch v.Count {
				case len(channels):
					if onReady != nil {
						if err := onReady(); err != nil {
							quit <- err
							return
						}
					}
				case 0:
					// 所有频道都解除订阅
					quit <- nil
					return
				}

			}
		}
	}()

	// 健康检查
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := psc.Ping(""); err != nil {
				quit <- err
				return
			}
		}
	}()

	// 自动重新订阅
	if num2 == -1 || num2 > 0 {
		go func() {
			<-quit
			time.Sleep(time.Second * 2)
			psc.Close()

			if num2 == -1 {
				_ = c.Subscribe(onReady, onMessage, num1, num2, channels...)
			} else {
				_ = c.Subscribe(onReady, onMessage, num1, num2-1, channels...)
			}
		}()
	}

	return nil
}
