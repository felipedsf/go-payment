package ledger

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/felipedsf/go-payment/ledger/internal/ledger"
	"log"
	"sync"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "root"
	dbName   = "ledger"
	topic    = "ledger"
)

var (
	db *sql.DB
	wg sync.WaitGroup
)

type LedgerMsg struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	Operation string `json:"operation"`
	Date      string `json:"date"`
}

func main() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPass, dbName)

	db, err = sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("Error closing db: %s", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

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
	var ledgerMsg LedgerMsg

	err := json.Unmarshal(msg.Value, &ledgerMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ledger.Insert(db, ledgerMsg.OrderID, ledgerMsg.UserID, ledgerMsg.Amount, ledgerMsg.Operation, ledgerMsg.Date)
	if err != nil {
		fmt.Println(err)
		return
	}
}
