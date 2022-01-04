package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/olivere/elastic/v7"
)

const MaxHandle = 1

var esPool *sync.Pool

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
	// client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
	// if err != nil {
	// 	log.Panic(err)
	// }
	client := esPool.Get().(*elastic.Client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := client.Index().Index("nsq-log").Type("log").BodyString(string(msg.Body)).Do(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(msg.Body))
	esPool.Put(client)
	return nil
}

func initPool() {
	esPool = &sync.Pool{
		New: func() interface{} {
			fmt.Printf("createing a new es client.\n")
			client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"), elastic.SetSniff(false))
			if err != nil {
				panic(err.Error())
			}
			return client
		},
	}
}

func main() {
	initPool()

	for i := 0; i < MaxHandle; i++ {
		go initConsumer(&LogConsumer{
			Topic:   "logrus-topic",
			Channel: "logrus-topic",
			Address: "localhost:4150",
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
