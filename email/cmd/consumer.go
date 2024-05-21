package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/felipedsf/go-payment/email/internal/email"
	"log"
	"sync"
)

const (
	topic = "email"
)

var wg *sync.WaitGroup

type EmailMsg struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}

func main() {
	done := make(chan struct{})

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, sarama.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		close(done)
		if err := consumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	for _, partition := range partitions {
		pConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = pConsumer.Close(); err != nil {
				log.Println(err)
			}
		}()

		wg.Add(1)
		go awaitMessages(done, pConsumer, partition)
	}
	wg.Wait()
}

func awaitMessages(done chan struct{}, partitionConsumer sarama.PartitionConsumer, partition int32) {
	defer wg.Done()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Partition %d - Received message: %s\n", partition, string(msg.Value))
			handleMessage(msg)
		case <-done:
			fmt.Printf("Received done signal. exiting...\n")
			return
		}
	}
}

func handleMessage(msg *sarama.ConsumerMessage) {
	var emailMsg EmailMsg

	err := json.Unmarshal(msg.Value, &emailMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = email.Send(emailMsg.OrderID, emailMsg.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}
}
