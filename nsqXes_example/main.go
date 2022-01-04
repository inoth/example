package main

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

const MaxHandle = 1

type RegisteredConsumer interface {
	GetTopic() string
	GetChannel() string
	GetAddress() string
}

type LogConsumer struct {
	Topic   string
	Channel string
	Address string
}

func (c *LogConsumer) GetTopic() string   { return c.Topic }
func (c *LogConsumer) GetChannel() string { return c.Channel }
func (c *LogConsumer) GetAddress() string { return c.Address }

func (c *LogConsumer) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("%v\n", string(msg.Body))

	return nil
}

func main() {
	for i := 0; i < MaxHandle; i++ {
		go initConsumer(&LogConsumer{
			Topic:   "logrus-topic",
			Channel: "logrus-topic",
			Address: "nsq host",
		})
	}
	select {}
}

func initConsumer(config RegisteredConsumer) {
	cfg := nsq.NewConfig()
	c, err := nsq.NewConsumer(config.GetTopic(), config.GetChannel(), cfg)
	if err != nil {
		panic(err)
	}
	c.AddHandler(config.(nsq.Handler))

	if err := c.ConnectToNSQD(config.GetAddress()); err != nil {
		panic(err)
	}
}
