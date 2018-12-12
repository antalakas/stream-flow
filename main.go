package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/quipo/statsd"
	"gitlab.crypteianetworks.prv/stream-flow/flow"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)
//lab.netflowdev2.flow.received
func main() {

	broker := "kafka1:9092"
	group := "group1"
	topics := []string{"netflow"}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	stats := getStatsClient("10.0.2.77:8125", "lab.netflowdev2.")

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

	counterThreshold := 10000.0
	counter := 1

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
			case *kafka.Stats:
				fmt.Println("==============================================")
				fmt.Printf("\n\nStats: %v\n\n", e)
				fmt.Println("==============================================")
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()

			case *kafka.Message:

				go analyzeEvent(c, e, stats)

				if math.Mod(float64(counter), counterThreshold) == 0 {
					counter = 1
					fmt.Println("10000 incoming events...")
				} else {
					counter = counter + 1
				}

			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
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

func getStatsClient(address string, prefix string) *statsd.StatsdBuffer {
	statsdClient := statsd.NewStatsdClient(address, prefix)
	err := statsdClient.CreateSocket()

	if nil != err {
		log.Println(err)
		//os.Exit(1)
	}

	interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	stats := statsd.NewStatsdBuffer(interval, statsdClient)
	//defer stats.Close()

	return stats
}

func analyzeEvent(
	c *kafka.Consumer,
	m *kafka.Message,
	stats *statsd.StatsdBuffer) {

	if m.Headers != nil {
		fmt.Printf("%% Headers: %v\n", m.Headers)
	}

	var netflowEvent flow.NetflowEvent

	messageContext := fmt.Sprintf("%s/%d/%d\t%s\t%s\n",
		m.TopicPartition.Topic, m.TopicPartition.Partition,
		m.TopicPartition.Offset, m.Key, m.Value)

	if err := json.Unmarshal(m.Value, &netflowEvent); err != nil {
		fmt.Println(messageContext)
		fmt.Println(err)
		return
	}

	//fmt.Println(messageContext)

	_, err := c.CommitMessage(m)

	//fmt.Println("Sending statistics")

	if err == nil {
		stats.Incr("flow.received", 1)
	}
}