package mq

import (
	"g"
	"github.com/streadway/amqp"
	//"log"
	"time"
)

var ConnPool *Pool

func InitConnPool() {

	ConnPool = &Pool{
		MaxIdle:     g.Config().Amqp.MaxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (*amqp.Connection, error) {
			var c *amqp.Connection
			var err error
			c, err = amqp.Dial(g.Config().Amqp.Addr)
			if err == nil {
				return c, nil
			}

			c, err = amqp.Dial(g.Config().Amqp.Addr1)
			if err == nil {
				return c, nil
			}

			c, err = amqp.Dial(g.Config().Amqp.Addr2)
			if err == nil {
				return c, nil
			}

			return nil, err
		},
		TestOnBorrow: nil,
	}
}
