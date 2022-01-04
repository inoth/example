package main

import (
	"fmt"
	"logrusXnsq_example/hook"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

func main() {
	config := nsq.NewConfig()
	client, err := nsq.NewProducer("nsq host", config)
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
	defer client.Stop()
	hook, err := hook.NewAsyncNsqHook(client, "logrus-topic", logrus.InfoLevel)

	logrus.AddHook(hook)

	logrus.Info("output to nsq")

	time.Sleep(time.Second * 5)
}
