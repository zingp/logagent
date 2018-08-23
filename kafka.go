package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var kafkaSend = &KafkaSend{}

type Message struct {
	line  string
	topic string
}

type KafkaSend struct {
	client   sarama.SyncProducer
	lineChan chan *Message
}

func initKafka(kafkaAddr string, threadNum int) (err error) {
	kafkaSend, err = NewKafkaSend(kafkaAddr, threadNum)
	return
}

// NewKafkaSend is 
func NewKafkaSend(kafkaAddr string, threadNum int) (kafka *KafkaSend, err error) {
	kafka = &KafkaSend{
		lineChan: make(chan *Message, 10000),
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // wait kafka ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // random partition
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init kafka client err: %v", err)
		return
	}
	kafka.client = client

	for i := 0; i < threadNum; i++ {
		fmt.Println("start to send kfk")
		waitGroup.Add(1)
		go kafka.sendMsgToKfk()
	}
	return
}

func (k *KafkaSend) sendMsgToKfk() {
	defer waitGroup.Done()

	for v := range k.lineChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.topic
		msg.Value = sarama.StringEncoder(v.line)

		_, _, err := k.client.SendMessage(msg)
		if err != nil {
			logs.Error("send massage to kafka error: %v", err)
			return
		}
	}
}

func (k *KafkaSend) addMessage(line string, topic string) (err error) {
	k.lineChan <- &Message{line: line, topic: topic}
	return
}
