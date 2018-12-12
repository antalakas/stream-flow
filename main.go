package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.crypteianetworks.prv/stream-flow/flow"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	broker := "kafka1:9092"
	group := "group1"
	topics := []string{"netflow"}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"group.id":           group,
		"session.timeout.ms": 6000,
		"auto.offset.reset":  "earliest"})

	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create consumer: %s\n", err))
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				//fmt.Printf("%% Message on %s:\n%s\n",
				//	e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}

				var netflowEvent flow.NetflowEvent

				messageContext := fmt.Sprintf("%s/%d/%d\t%s\t%s\n",
					e.TopicPartition.Topic, e.TopicPartition.Partition,
					e.TopicPartition.Offset, e.Key, e.Value)

				if err := json.Unmarshal(e.Value, &netflowEvent); err != nil {
					fmt.Println(messageContext)
					fmt.Println(err)
					continue
				}

			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Println(fmt.Sprintf("%% Error: %v\n", e))
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	err = c.Close()
}